package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/genudine/saerro-go/translators"
	"github.com/genudine/saerro-go/types"
	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

func getEventHandlerTestShim(t *testing.T) (EventHandler, *sql.DB) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("test shim: sqlite open failed, %v", err)
	}

	return EventHandler{
		Ingest: &Ingest{
			DB: db,
		},
	}, db
}

func TestHandleDeath(t *testing.T) {
	eh, db := getEventHandlerTestShim(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	event := types.ESSEvent{
		EventName: "Death",
		WorldID:   17,
		ZoneID:    2,

		CharacterID: "DollNC",
		LoadoutID:   "3",
		TeamID:      types.NC,

		AttackerCharacterID: "Lyyti",
		AttackerLoadoutID:   "3",
		AttackerTeamID:      types.TR,
	}

	eh.HandleDeath(ctx, event)

	type player struct {
		CharacterID string `json:"character_id"`
		ClassName   string `json:"class_name"`
	}

	var player1 player
	err := db.QueryRowContext(ctx, "SELECT * FROM players WHERE character_id = ?", event.CharacterID).Scan(&player1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, event.CharacterID, player1.CharacterID)
	assert.Equal(t, translators.ClassFromLoadout(event.LoadoutID), player1.ClassName)

	var player2 player
	err = db.QueryRowContext(ctx, "SELECT * FROM players WHERE character_id = ?", event.AttackerCharacterID).Scan(&player2)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, event.AttackerCharacterID, player2.CharacterID)
	assert.Equal(t, translators.ClassFromLoadout(event.AttackerLoadoutID), player2.ClassName)
}
