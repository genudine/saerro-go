package store_test

import (
	"testing"
	"time"

	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/types"

	"github.com/stretchr/testify/assert"
)

func TestVehicleMigrationClean(t *testing.T) {
	ctx, db := mkEnv(t)

	ps := store.NewVehicleStore(db)

	isMigrated := ps.IsMigrated(ctx)
	assert.True(t, isMigrated)
}

func TestVehicleMultipleStartups(t *testing.T) {
	ctx, db := mkEnv(t)

	ps1 := store.NewVehicleStore(db)
	ps2 := store.NewVehicleStore(db)

	isMigrated := ps2.IsMigrated(ctx)
	assert.True(t, isMigrated)

	assert.True(t, ps1.RanMigration)
	assert.False(t, ps2.RanMigration)
}

func TestVehicleMigrationRerun(t *testing.T) {
	ctx, db := mkEnv(t)

	ps := store.NewVehicleStore(db)
	ps.RunMigration(ctx, true)

	isMigrated := ps.IsMigrated(ctx)
	assert.True(t, isMigrated)
}

func TestVehicleInsertGetOne(t *testing.T) {
	ctx, db := mkEnv(t)
	ps := store.NewVehicleStore(db)

	player := &types.Vehicle{
		CharacterID: "41666",
		WorldID:     17,
		FactionID:   types.TR,
		ZoneID:      4,
		VehicleName: "harasser",
	}

	err := ps.Insert(ctx, player)
	assert.NoError(t, err, "Insert failed")

	p1, err := ps.GetOne(ctx, "41666")
	assert.NoError(t, err, "GetOne failed")
	assert.Equal(t, "harasser", p1.VehicleName)

	time.Sleep(time.Second * 1)

	player.VehicleName = "vanguard"
	player.LastUpdated = time.Time{}
	err = ps.Insert(ctx, player)
	assert.NoError(t, err, "Insert failed")

	p2, err := ps.GetOne(ctx, "41666")
	assert.NoError(t, err, "GetOne failed")
	assert.Equal(t, "vanguard", p2.VehicleName)
	assert.NotEqual(t, p1.LastUpdated, p2.LastUpdated, "time did not update as expected")
}

func TestVehiclePrune(t *testing.T) {
	ctx, db := mkEnv(t)
	ps := store.NewVehicleStore(db)

	prunedPlayer := &types.Vehicle{
		CharacterID: "20155",
		WorldID:     17,
		FactionID:   types.NC,
		ZoneID:      4,
		VehicleName: "harasser",

		LastUpdated: time.Now().Add(-time.Minute * 20),
	}
	survivingPlayer := &types.Vehicle{
		CharacterID: "41666",
		WorldID:     17,
		FactionID:   types.TR,
		ZoneID:      4,
		VehicleName: "harasser",
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
