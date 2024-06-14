package main

import (
	"context"
	"log"

	"github.com/genudine/saerro-go/types"
)

type EventHandler struct {
	Ingest *Ingest
}

func (eh *EventHandler) HandleEvent(ctx context.Context, event types.ESSEvent) {
	if event.EventName == "" {
		log.Println("invalid event; dropping")
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
		log.Println("got pop event")
		pe := PopEventFromESSEvent(event, false)
		eh.Ingest.TrackPop(ctx, pe)
	}

	if event.AttackerCharacterID != "" && event.AttackerCharacterID != "0" && event.AttackerTeamID != 0 {
		log.Println("got attacker pop event")
		pe := PopEventFromESSEvent(event, true)
		eh.Ingest.TrackPop(ctx, pe)
	}
}

func (eh *EventHandler) HandleExperience(ctx context.Context, event types.ESSEvent) {
	// Detect specific vehicles via related experience IDs
	vehicleID := ""
	switch event.ExperienceID {
	case 201: // Galaxy Spawn Bonus
		vehicleID = "11"
		break
	case 233: // Sunderer Spawn Bonus
		vehicleID = "2"
		break
	case 674: // ANT stuff
	case 675:
		vehicleID = "160"
		break
	}

	event.VehicleID = vehicleID
	pe := PopEventFromESSEvent(event, false)
	eh.Ingest.TrackPop(ctx, pe)
}

func (eh *EventHandler) HandleAnalytics(ctx context.Context, event types.ESSEvent) {

}
