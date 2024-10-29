package ingest_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genudine/saerro-go/cmd/ws/ingest"
	"github.com/genudine/saerro-go/store/storemock"
	"github.com/genudine/saerro-go/translators"
	"github.com/genudine/saerro-go/types"
)

func mkIngest(t *testing.T) (context.Context, *ingest.Ingest, *storemock.MockPlayerStore, *storemock.MockVehicleStore) {
	t.Helper()

	ps := new(storemock.MockPlayerStore)
	vs := new(storemock.MockVehicleStore)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	t.Cleanup(cancel)

	i := &ingest.Ingest{
		PlayerStore:  ps,
		VehicleStore: vs,
	}

	return ctx, i, ps, vs
}

func TestTrackPopHappyPath(t *testing.T) {
	ctx, i, ps, _ := mkIngest(t)

	// Combat Medic on Emerald
	event := types.PopEvent{
		WorldID:     17,
		ZoneID:      4,
		TeamID:      types.TR,
		LoadoutID:   4,
		ClassName:   translators.CombatMedic,
		CharacterID: "aaaa",
	}

	eventPlayer := event.ToPlayer()

	ps.On("Insert", ctx, eventPlayer).Return(nil).Once()

	i.TrackPop(ctx, event)
}

func TestTrackPopFixup(t *testing.T) {
	ctx, i, ps, _ := mkIngest(t)

	event := types.PopEvent{
		WorldID:     17,
		ZoneID:      4,
		TeamID:      0,
		ClassName:   "unknown",
		CharacterID: "bbbb",
	}
	pastEventPlayer := event.ToPlayer()
	pastEventPlayer.ClassName = "light_assault"
	pastEventPlayer.FactionID = types.VS

	ps.On("GetOne", ctx, event.CharacterID).Return(pastEventPlayer, nil).Once()
	ps.On("Insert", ctx, pastEventPlayer).Return(nil).Once()

	i.TrackPop(ctx, event)
}

func TestTrackPopFixupFailed(t *testing.T) {
	ctx, i, ps, _ := mkIngest(t)

	event := types.PopEvent{
		WorldID:     17,
		ZoneID:      4,
		TeamID:      0,
		ClassName:   "unknown",
		CharacterID: "bbbb",
	}

	ps.On("GetOne", ctx, event.CharacterID).Return(nil, errors.New("ingest fixup failed")).Once()

	i.TrackPop(ctx, event)
}
