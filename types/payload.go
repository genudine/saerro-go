package types

type ESSData struct {
	Payload ESSEvent
}

type ESSEvent struct {
	EventName string `json:"event_name"`
	WorldID   uint16 `json:"world_id"`
	ZoneID    uint32 `json:"zone_id"`

	CharacterID string  `json:"character_id"`
	LoadoutID   string  `json:"loadout_id"`
	VehicleID   string  `json:"vehicle_id"`
	TeamID      Faction `json:"team_id"`

	AttackerCharacterID string  `json:"attacker_character_id"`
	AttackerLoadoutID   string  `json:"attacker_loadout_id"`
	AttackerVehicleID   string  `json:"attacker_vehicle_id"`
	AttackerTeamID      Faction `json:"attacker_team_id"`

	ExperienceID uint32
}
