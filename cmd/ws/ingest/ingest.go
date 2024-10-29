package ingest

import (
	"context"
	"fmt"
	"log"

	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/types"
)

type IIngest interface {
	TrackPop(context.Context, types.PopEvent)
}

type Ingest struct {
	PlayerStore  store.IPlayerStore
	VehicleStore store.IVehicleStore
}

func (i *Ingest) TrackPop(ctx context.Context, event types.PopEvent) {
	player := event.ToPlayer()

	err := i.fixupPlayer(ctx, player)
	if err != nil {
		log.Println("ingest: player fixup failed, dropping event", err)
		return
	}

	err = i.PlayerStore.Insert(ctx, player)
	if err != nil {
		log.Println("TrackPop Insert failed", err)
	}
}

func (i *Ingest) fixupPlayer(ctx context.Context, player *types.Player) error {
	if player.ClassName != "unknown" && player.FactionID != 0 {
		// all fixups are done
		return nil
	}

	storedPlayer, err := i.PlayerStore.GetOne(ctx, player.CharacterID)
	if err != nil {
		return fmt.Errorf("ingest: fixupPlayer: fetching player %s failed: %w", player.CharacterID, err)
	}

	// probably VehicleDestroy
	if player.ClassName == "unknown" {
		// TODO: maybe get this from census, profile_id
		player.ClassName = storedPlayer.ClassName
	}

	// probably PS4
	if player.FactionID == 0 {
		// TODO: get this from census
		player.FactionID = storedPlayer.FactionID
	}

	return nil
}
