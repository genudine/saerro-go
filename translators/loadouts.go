package translators

type Class string

const (
	Infiltrator  Class = "infiltrator"
	LightAssault Class = "light_assault"
	CombatMedic  Class = "combat_medic"
	Engineer     Class = "engineer"
	HeavyAssault Class = "heavy_assault"
	MAX          Class = "max"
)

var (
	LoadoutMap = map[uint16]Class{
		1:  Infiltrator,
		8:  Infiltrator,
		15: Infiltrator,
		28: Infiltrator,

		3:  LightAssault,
		10: LightAssault,
		17: LightAssault,
		29: LightAssault,

		4:  CombatMedic,
		11: CombatMedic,
		18: CombatMedic,
		30: CombatMedic,

		5:  Engineer,
		12: Engineer,
		19: Engineer,
		31: Engineer,

		6:  HeavyAssault,
		13: HeavyAssault,
		20: HeavyAssault,
		32: HeavyAssault,

		7:  MAX,
		14: MAX,
		21: MAX,
		45: MAX,
	}
)

func ClassFromLoadout(loadoutID uint16) Class {
	c, ok := LoadoutMap[loadoutID]

	if !ok {
		return "unknown"
	}

	return c
}
