package translators

type Vehicle string

const (
	Flash     Vehicle = "flash"
	Sunderer  Vehicle = "sunderer"
	Lightning Vehicle = "lightning"
	Magrider  Vehicle = "magrider"
	Vanguard  Vehicle = "vanguard"
	Prowler   Vehicle = "prowler"
	Scythe    Vehicle = "scythe"
	Reaver    Vehicle = "reaver"
	Mosquito  Vehicle = "mosquito"
	Liberator Vehicle = "liberator"
	Galaxy    Vehicle = "galaxy"
	Harasser  Vehicle = "harasser"
	Valkyrie  Vehicle = "valkyrie"
	Ant       Vehicle = "ant"
	Dervish   Vehicle = "dervish"
	Chimera   Vehicle = "chimera"
	Javelin   Vehicle = "javelin"
	Corsair   Vehicle = "corsair"
)

func VehicleNameFromID(id string) Vehicle {
	v, ok := VehicleMap[id]
	if !ok {
		return "unknown"
	}

	return v
}
