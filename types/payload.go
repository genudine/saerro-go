package types

type ESSData struct {
	Payload ESSEvent
}

type ESSEvent struct {
	EventName string `json:"event_name"`
	WorldID   uint16 `json:"world_id,string"`
	ZoneID    uint32 `json:"zone_id,string"`

	CharacterID string `json:"character_id"`

	// On Death

	VehicleID           string  `json:"vehicle_id"`
	TeamID              Faction `json:"team_id,string"`
	CharacterLoadoutID  uint16  `json:"character_loadout_id,string"`
	AttackerCharacterID string  `json:"attacker_character_id"`
	AttackerLoadoutID   uint16  `json:"attacker_loadout_id,string"`
	AttackerVehicleID   string  `json:"attacker_vehicle_id"`
	AttackerTeamID      Faction `json:"attacker_team_id,string"`

	// On GainExperience

	ExperienceID uint32 `json:"experience_id,string"`
	LoadoutID    uint16 `json:"loadout_id,string"`
}
