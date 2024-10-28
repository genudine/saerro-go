package types

type Faction uint8

const (
	VS Faction = iota
	NC
	TR
	NSO
)

// func (f Faction) UnmarshalJSON(b []byte) error {
// 	switch b[0] {
// 	case '1':
// 		f = VS
// 	case '2':
// 		f = NC
// 	case '3':
// 		f = TR
// 	case '4':
// 		f = NSO
// 	}

// 	return nil
// }
