package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/genudine/saerro-go/types"

	"github.com/avast/retry-go"
)

type IVehicleStore interface {
	IsMigrated(context.Context) bool
	RunMigration(context.Context, bool)
	Insert(context.Context, *types.Vehicle) error
	GetOne(context.Context, string) (*types.Vehicle, error)
	Prune(context.Context) (int64, error)
}

type VehicleStore struct {
	DB *sql.DB

	// Test introspection for if migrations ran during this PlayerStore init
	RanMigration bool
}

func NewVehicleStore(db *sql.DB) *VehicleStore {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ps := &VehicleStore{
		DB: db,
	}

	ps.RunMigration(ctx, false)

	return ps
}

func (ps *VehicleStore) IsMigrated(ctx context.Context) bool {
	_, err := ps.DB.QueryContext(ctx, `SELECT count(1) FROM vehicles LIMIT 1;`)
	if err != nil {
		log.Printf("IsMigrated check failed: %v", err)
		return false
	}

	return true
}

func (ps *VehicleStore) RunMigration(ctx context.Context, force bool) {
	if !force && ps.IsMigrated(ctx) {
		return
	}

	log.Println("(Re)-creating vehicles table")
	ps.DB.ExecContext(ctx, `
		DROP TABLE IF EXISTS vehicles;

		CREATE TABLE vehicles (
			character_id TEXT NOT NULL PRIMARY KEY,
			last_updated TIMESTAMP NOT NULL,
			world_id INT NOT NULL,
			faction_id INT NOT NULL,
			zone_id INT NOT NULL,
			vehicle_name TEXT NOT NULL
		);

		-- TODO: Add indexes?
	`)
	log.Println("Done, vehicles table is initialized.")
	ps.RanMigration = true
}

// Insert a player into the store.
// For testing, when LastUpdated is not "zero", the provided timestamp will be carried into the store.
func (ps *VehicleStore) Insert(ctx context.Context, vehicle *types.Vehicle) error {
	if vehicle.LastUpdated.IsZero() {
		vehicle.LastUpdated = time.Now()
	}

	err := retry.Do(func() error {
		_, err := ps.DB.ExecContext(ctx,
			`
			INSERT INTO vehicles (
				last_updated, character_id, world_id, faction_id, zone_id, vehicle_name
			) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			ON CONFLICT (character_id) DO UPDATE SET 
				last_updated = EXCLUDED.last_updated,
				world_id = EXCLUDED.world_id,
				faction_id = EXCLUDED.faction_id,
				zone_id = EXCLUDED.zone_id,
				vehicle_name = EXCLUDED.vehicle_name
			`,
			vehicle.LastUpdated,
			vehicle.CharacterID,
			vehicle.WorldID,
			vehicle.FactionID,
			vehicle.ZoneID,
			vehicle.VehicleName,
		)

		return err
	}, retry.Attempts(2))

	return err
}

// GetOne player from the store.
func (ps *VehicleStore) GetOne(ctx context.Context, id string) (*types.Vehicle, error) {
	row := ps.DB.QueryRowContext(ctx, `
		SELECT 
			last_updated,
			character_id, 
			world_id, 
			faction_id, 
			zone_id, 
			vehicle_name
		FROM vehicles WHERE character_id = $1
	`, id)

	vehicle := &types.Vehicle{}
	var timestamp string

	err := row.Scan(
		&timestamp,
		&vehicle.CharacterID,
		&vehicle.WorldID,
		&vehicle.FactionID,
		&vehicle.ZoneID,
		&vehicle.VehicleName,
	)
	if err != nil {
		return nil, err
	}

	vehicle.LastUpdated, err = time.Parse(time.RFC3339, timestamp)

	return vehicle, err
}

func (ps *VehicleStore) Prune(ctx context.Context) (int64, error) {
	log.Println("pruning VehicleStore")

	// Avoid using sql idioms here for portability
	// SQLite and PgSQL do now() differently, we don't need to at all.
	res, err := ps.DB.ExecContext(ctx, `
		DELETE FROM vehicles WHERE last_updated < $1;
	`, time.Now().Add(-time.Minute*15))
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
