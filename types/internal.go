package types

import "time"

type Player struct {
	CharacterID string    `json:"character_id"`
	LastUpdated time.Time `json:"last_updated"`
	WorldID     uint16    `json:"world_id"`
	FactionID   Faction   `json:"faction_id"`
	ZoneID      uint32    `json:"zone_id"`
	ClassName   string    `json:"class_name"`
}

type Vehicle struct {
	CharacterID string    `json:"character_id"`
	LastUpdated time.Time `json:"last_updated"`
	WorldID     uint16    `json:"world_id"`
	FactionID   Faction   `json:"faction_id"`
	ZoneID      uint32    `json:"zone_id"`
	VehicleName string    `json:"vehicle_name"`
}

type AnalyticEvent struct {
	Time      time.Time `json:"time"`
	WorldID   uint16    `json:"world_id"`
	EventName string    `json:"event_name"`
}
