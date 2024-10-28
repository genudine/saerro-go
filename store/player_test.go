package store_test

import (
	"testing"
	"time"

	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/types"

	"github.com/stretchr/testify/assert"
)

func TestPlayerMigrationClean(t *testing.T) {
	ctx, db := mkEnv(t)

	ps := store.NewPlayerStore(db)

	isMigrated := ps.IsMigrated(ctx)
	assert.True(t, isMigrated)
}

func TestPlayerMultipleStartups(t *testing.T) {
	ctx, db := mkEnv(t)

	ps1 := store.NewPlayerStore(db)
	ps2 := store.NewPlayerStore(db)

	isMigrated := ps2.IsMigrated(ctx)
	assert.True(t, isMigrated)

	assert.True(t, ps1.RanMigration)
	assert.False(t, ps2.RanMigration)
}

func TestPlayerMigrationRerun(t *testing.T) {
	ctx, db := mkEnv(t)

	ps := store.NewPlayerStore(db)
	ps.RunMigration(ctx, true)

	isMigrated := ps.IsMigrated(ctx)
	assert.True(t, isMigrated)
}

func TestPlayerInsertGetOne(t *testing.T) {
	ctx, db := mkEnv(t)
	ps := store.NewPlayerStore(db)

	player := &types.Player{
		CharacterID: "41666",
		WorldID:     17,
		FactionID:   types.TR,
		ZoneID:      4,
		ClassName:   "light_assault",
	}

	err := ps.Insert(ctx, player)
	assert.NoError(t, err, "Insert failed")

	p1, err := ps.GetOne(ctx, "41666")
	assert.NoError(t, err, "GetOne failed")
	assert.Equal(t, "light_assault", p1.ClassName)

	time.Sleep(time.Second * 1)

	player.ClassName = "combat_medic"
	player.LastUpdated = time.Time{}
	err = ps.Insert(ctx, player)
	assert.NoError(t, err, "Insert failed")

	p2, err := ps.GetOne(ctx, "41666")
	assert.NoError(t, err, "GetOne failed")
	assert.Equal(t, "combat_medic", p2.ClassName)
	assert.NotEqual(t, p1.LastUpdated, p2.LastUpdated, "time did not update as expected")
}

func TestPlayerPrune(t *testing.T) {
	ctx, db := mkEnv(t)
	ps := store.NewPlayerStore(db)

	prunedPlayer := &types.Player{
		CharacterID: "20155",
		WorldID:     17,
		FactionID:   types.NC,
		ZoneID:      4,
		ClassName:   "light_assault",
		LastUpdated: time.Now().Add(-time.Minute * 20),
	}
	survivingPlayer := &types.Player{
		CharacterID: "41666",
		WorldID:     17,
		FactionID:   types.TR,
		ZoneID:      4,
		ClassName:   "light_assault",
		LastUpdated: time.Now().Add(-time.Minute * 5),
	}

	err := ps.Insert(ctx, prunedPlayer)
	assert.NoError(t, err, "Insert prunedPlayer failed")

	err = ps.Insert(ctx, survivingPlayer)
	assert.NoError(t, err, "Insert survivingPlayer failed")

	removed, err := ps.Prune(ctx)
	assert.NoError(t, err, "Prune failed")
	assert.Equal(t, int64(1), removed, "Prune count incorrect")

	_, err = ps.GetOne(ctx, prunedPlayer.CharacterID)
	assert.Error(t, err, "GetOne prunedPlayer failed, as expected.")

	_, err = ps.GetOne(ctx, survivingPlayer.CharacterID)
	assert.NoError(t, err, "GetOne survivingPlayer failed")
}
