package eventhandler

import (
	"context"
	"database/sql"

	"github.com/genudine/saerro-go/cmd/ws/ingest"
	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/types"
)

type EventHandler struct {
	Ingest *ingest.Ingest
}

func NewEventHandler(db *sql.DB) EventHandler {
	return EventHandler{
		Ingest: &ingest.Ingest{
			PlayerStore: store.NewPlayerStore(db),
		},
	}
}

func (eh *EventHandler) HandleEvent(ctx context.Context, event types.ESSEvent) {
	if event.EventName == "" {
		// log.Println("invalid event; dropping")
		return
	}

	if event.EventName == "Death" || event.EventName == "VehicleDestroy" {
		go eh.HandleDeath(ctx, event)
	} else if event.EventName == "GainExperience" {
		go eh.HandleExperience(ctx, event)
	}

	go eh.HandleAnalytics(ctx, event)
}

func (eh *EventHandler) HandleDeath(ctx context.Context, event types.ESSEvent) {
	if event.CharacterID != "" && event.CharacterID != "0" {
		// log.Println("got pop event")
		pe := types.PopEventFromESSEvent(event, false)
		eh.Ingest.TrackPop(ctx, pe)
	}

	if event.AttackerCharacterID != "" && event.AttackerCharacterID != "0" && event.AttackerTeamID != 0 {
		pe := types.PopEventFromESSEvent(event, true)
		// fmt.Println("got attacker pop event", event)
		eh.Ingest.TrackPop(ctx, pe)
	}
}

func (eh *EventHandler) HandleExperience(ctx context.Context, event types.ESSEvent) {
	// Detect specific vehicles via related experience IDs
	vehicleID := ""
	switch event.ExperienceID {
	case 201: // Galaxy Spawn Bonus
		vehicleID = "11"
	case 233: // Sunderer Spawn Bonus
		vehicleID = "2"
	case 674:
		fallthrough // ANT stuff
	case 675:
		vehicleID = "160"
	}

	event.VehicleID = vehicleID
	pe := types.PopEventFromESSEvent(event, false)
	eh.Ingest.TrackPop(ctx, pe)
}

func (eh *EventHandler) HandleAnalytics(ctx context.Context, event types.ESSEvent) {

}
