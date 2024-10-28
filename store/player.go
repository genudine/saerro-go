package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/genudine/saerro-go/types"
	"github.com/genudine/saerro-go/util"

	"github.com/avast/retry-go"
)

type IPlayerStore interface {
	IsMigrated(context.Context) bool
	RunMigration(context.Context, bool)
	Insert(context.Context, *types.Player) error
	GetOne(context.Context, string) (*types.Player, error)
	Prune(context.Context) (int64, error)
}

type PlayerStore struct {
	DB *sql.DB

	// Test introspection for if migrations ran during this PlayerStore init
	RanMigration bool
}

func NewPlayerStore(db *sql.DB) *PlayerStore {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ps := &PlayerStore{
		DB: db,
	}

	ps.RunMigration(ctx, false)

	return ps
}

func (ps *PlayerStore) IsMigrated(ctx context.Context) bool {
	_, err := ps.DB.QueryContext(ctx, `SELECT count(1) FROM players LIMIT 1;`)
	if err != nil {
		log.Printf("IsMigrated check failed: %v", err)
		return false
	}

	return true
}

func (ps *PlayerStore) RunMigration(ctx context.Context, force bool) {
	if !force && ps.IsMigrated(ctx) {
		return
	}

	log.Println("(Re)-creating players table")
	ps.DB.ExecContext(ctx, `
		DROP TABLE IF EXISTS players;

		CREATE TABLE players (
			character_id TEXT NOT NULL PRIMARY KEY,
			last_updated TIMESTAMP NOT NULL,
			world_id INT NOT NULL,
			faction_id INT NOT NULL,
			zone_id INT NOT NULL,
			class_name TEXT NOT NULL
		);

		-- TODO: Add indexes?
	`)
	log.Println("Done, players table is initialized.")
	ps.RanMigration = true
}

// Insert a player into the store.
// For testing, when LastUpdated is not "zero", the provided timestamp will be carried into the store.
func (ps *PlayerStore) Insert(ctx context.Context, player *types.Player) error {
	if player.LastUpdated.IsZero() {
		player.LastUpdated = time.Now()
	}

	err := retry.Do(func() error {
		_, err := ps.DB.ExecContext(ctx,
			`
			INSERT INTO players (
				last_updated, character_id, world_id, faction_id, zone_id, class_name
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (character_id) DO UPDATE SET
				last_updated = EXCLUDED.last_updated,
				world_id = EXCLUDED.world_id,
				faction_id = EXCLUDED.faction_id,
				zone_id = EXCLUDED.zone_id,
				class_name = EXCLUDED.class_name;
			`,
			util.TimeToString(player.LastUpdated),
			player.CharacterID,
			player.WorldID,
			player.FactionID,
			player.ZoneID,
			player.ClassName,
		)

		return err
	}, retry.Attempts(2))

	return err
}

// GetOne player from the store.
func (ps *PlayerStore) GetOne(ctx context.Context, id string) (*types.Player, error) {
	row := ps.DB.QueryRowContext(ctx, `
		SELECT
			last_updated,
			character_id,
			world_id,
			faction_id,
			zone_id,
			class_name
		FROM players WHERE character_id = $1
	`, id)

	player := &types.Player{}
	var timestamp string

	err := row.Scan(
		&timestamp,
		&player.CharacterID,
		&player.WorldID,
		&player.FactionID,
		&player.ZoneID,
		&player.ClassName,
	)
	if err != nil {
		return nil, err
	}

	player.LastUpdated, err = time.Parse(time.RFC3339, timestamp)

	return player, err
}

func (ps *PlayerStore) Prune(ctx context.Context) (int64, error) {
	log.Println("pruning PlayerStore")

	// Avoid using sql idioms here for portability
	// SQLite and PgSQL do now() differently, we don't need to at all.
	res, err := ps.DB.ExecContext(ctx,
		`
			DELETE FROM players WHERE last_updated < $1;
		`,
		util.TimeToString(time.Now().Add(-time.Minute*15)),
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
