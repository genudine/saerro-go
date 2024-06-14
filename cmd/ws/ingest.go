package main

import (
	"context"
	"database/sql"
)

type Ingest struct {
	DB *sql.DB
}

func (i *Ingest) TrackPop(ctx context.Context, event PopEvent) {

}
