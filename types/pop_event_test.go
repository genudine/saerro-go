package types_test

import (
	"testing"

	"github.com/genudine/saerro-go/types"
	"github.com/stretchr/testify/assert"
)

var (
	deathEvent = types.ESSEvent{
		EventName: "Death",
		WorldID:   17,
		ZoneID:    4,

		CharacterID:        "LyytisDoll",
		TeamID:             types.NC,
		CharacterLoadoutID: 4,

		AttackerCharacterID: "Lyyti",
		AttackerTeamID:      types.TR,
		AttackerLoadoutID:   3,
	}

	expEvent = types.ESSEvent{
		EventName: "GainExperience",
		WorldID:   17,
		ZoneID:    4,

		CharacterID: "LyytisDoll",
		TeamID:      types.NC,
		LoadoutID:   3,

		ExperienceID: 55,
	}
)

func TestPopEventFromDeath(t *testing.T) {
	victim := types.PopEventFromESSEvent(deathEvent, false)
	assert.Equal(t, deathEvent.CharacterID, victim.CharacterID)
	assert.Equal(t, deathEvent.TeamID, victim.TeamID)
	assert.Equal(t, deathEvent.CharacterLoadoutID, victim.LoadoutID)
	assert.Equal(t, "combat_medic", string(victim.ClassName))

	attacker := types.PopEventFromESSEvent(deathEvent, true)
	assert.Equal(t, deathEvent.AttackerCharacterID, attacker.CharacterID)
	assert.Equal(t, deathEvent.AttackerTeamID, attacker.TeamID)
	assert.Equal(t, deathEvent.AttackerLoadoutID, attacker.LoadoutID)
	assert.Equal(t, "light_assault", string(attacker.ClassName))
}

func TestPopEventFromExperienceGain(t *testing.T) {
	pe := types.PopEventFromESSEvent(expEvent, false)
	assert.Equal(t, expEvent.CharacterID, pe.CharacterID)
	assert.Equal(t, expEvent.TeamID, pe.TeamID)
	assert.Equal(t, expEvent.LoadoutID, pe.LoadoutID)
	assert.Equal(t, "light_assault", string(pe.ClassName))
}

func TestPopEventFromVehicleDestroy(t *testing.T) {
	t.SkipNow()
}

func TestPopEventToPlayer(t *testing.T) {
	pe := types.PopEventFromESSEvent(deathEvent, false)
	player := pe.ToPlayer()
	assert.Equal(t, pe.CharacterID, player.CharacterID)
	assert.Equal(t, pe.TeamID, player.FactionID)
	assert.Equal(t, pe.ZoneID, player.ZoneID)
	assert.Equal(t, string(pe.ClassName), player.ClassName)
	assert.Equal(t, pe.WorldID, player.WorldID)
}
