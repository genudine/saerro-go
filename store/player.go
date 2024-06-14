package store

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type PlayerStore struct {
	DB *sql.DB
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

func (ps *PlayerStore) RunMigration(ctx context.Context, force bool) {
	if !force {
		// check if migrated first...

	}

	log.Println("(Re)-creating players table")
	ps.DB.ExecContext(ctx, `
		DROP TABLE IF EXISTS players;

		CREATE TABLE players (
			character_id TEXT NOT NULL PRIMARY KEY,
			last_updated TIMESTAMPTZ NOT NULL,
			world_id INT NOT NULL,
			faction_id INT NOT NULL,
			zone_id INT NOT NULL,
			class_name TEXT NOT NULL
		);
	`)

	log.Println("Done, players table is initialized.")
}
