// TODO: tests
package types

import (
	"github.com/genudine/saerro-go/translators"
)

type PopEvent struct {
	WorldID     uint16
	ZoneID      uint32
	CharacterID string
	LoadoutID   uint16
	TeamID      Faction
	VehicleID   string

	VehicleName translators.Vehicle
	ClassName   translators.Class
}

func PopEventFromESSEvent(event ESSEvent, attacker bool) PopEvent {
	pe := PopEvent{
		WorldID: event.WorldID,
		ZoneID:  event.ZoneID,
	}

	if !attacker {
		pe.CharacterID = event.CharacterID
		pe.LoadoutID = event.CharacterLoadoutID
		pe.TeamID = event.TeamID
		pe.VehicleID = event.VehicleID
	} else {
		pe.CharacterID = event.AttackerCharacterID
		pe.LoadoutID = event.AttackerLoadoutID
		pe.TeamID = event.AttackerTeamID
		pe.VehicleID = event.AttackerVehicleID
	}

	if pe.LoadoutID == 0 {
		pe.LoadoutID = event.LoadoutID
	}

	pe.ClassName = translators.ClassFromLoadout(pe.LoadoutID)
	pe.VehicleName = translators.VehicleNameFromID(pe.VehicleID)

	return pe
}

func (pe PopEvent) ToPlayer() *Player {
	return &Player{
		CharacterID: pe.CharacterID,
		ClassName:   string(pe.ClassName),
		FactionID:   pe.TeamID,
		ZoneID:      pe.ZoneID,
		WorldID:     pe.WorldID,
	}
}
