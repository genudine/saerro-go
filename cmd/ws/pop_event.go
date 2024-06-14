package main

import (
	"github.com/genudine/saerro-go/translators"
	"github.com/genudine/saerro-go/types"
)

type PopEvent struct {
	WorldID     uint16
	ZoneID      uint32
	CharacterID string
	LoadoutID   string
	TeamID      types.Faction
	VehicleID   string

	VehicleName translators.Vehicle
	ClassName   translators.Class
}

func PopEventFromESSEvent(event types.ESSEvent, attacker bool) PopEvent {
	pe := PopEvent{
		WorldID: event.WorldID,
		ZoneID:  event.ZoneID,
	}

	if !attacker {
		pe.CharacterID = event.CharacterID
		pe.LoadoutID = event.LoadoutID
		pe.TeamID = event.TeamID
		pe.VehicleID = event.VehicleID
	} else {
		pe.CharacterID = event.AttackerCharacterID
		pe.LoadoutID = event.AttackerLoadoutID
		pe.TeamID = event.AttackerTeamID
		pe.VehicleID = event.AttackerVehicleID
	}

	pe.ClassName = translators.ClassFromLoadout(pe.LoadoutID)
	pe.VehicleName = translators.VehicleNameFromID(pe.VehicleID)

	return pe
}
