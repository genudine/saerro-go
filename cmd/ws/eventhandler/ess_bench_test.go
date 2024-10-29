package eventhandler_test

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/genudine/saerro-go/cmd/ws/eventhandler"
	"github.com/genudine/saerro-go/translators"
	"github.com/genudine/saerro-go/types"
	"github.com/genudine/saerro-go/util"
)

const PreloadedCharacterCount = 45

var (
	characterStore = make([]string, PreloadedCharacterCount)
	db             *sql.DB
)

func BenchmarkESS(b *testing.B) {
	events, eh := prebench(b)

	for i := 0; i < b.N; i++ {
		event := events[i]

		eh.HandleEvent(context.Background(), event)
	}
}

func prebench(b *testing.B) ([]types.ESSEvent, eventhandler.EventHandler) {
	b.Helper()
	b.StopTimer()

	// Create our character pool
	for i := 0; i < PreloadedCharacterCount; i++ {
		characterStore[i] = mkRandomCharacterID()
	}

	// Create events
	events := []types.ESSEvent{}
	for i := 0; i < b.N; i++ {
		events = append(events, mkRandomEvent())
	}

	if db == nil {
		var err error
		db, err = util.GetDBConnection(os.Getenv("DB_ADDR"))
		if err != nil {
			b.Fatal(err)
		}
	}

	eh := eventhandler.NewEventHandler(db)

	b.ResetTimer()
	b.StartTimer()

	return events, eh
}

func mkRandomEvent() types.ESSEvent {
	w := rand.Intn(4)
	z := rand.Intn(7)

	switch rand.Intn(2) {
	case 0:
		return mkRandomDeathEvent(w, z)
	default:
		return mkRandomExpEvent(w, z)
	}
}

func mkRandomDeathEvent(world, zone int) types.ESSEvent {
	return types.ESSEvent{
		EventName: "Death",
		WorldID:   uint16(world),
		ZoneID:    uint32(zone),

		CharacterID:        getRandomCharacterID(),
		CharacterLoadoutID: mkRandomLoadout(),
		TeamID:             mkRandomFaction(),

		AttackerCharacterID: getRandomCharacterID(),
		AttackerLoadoutID:   mkRandomLoadout(),
		AttackerTeamID:      mkRandomFaction(),
	}
}

func mkRandomExpEvent(world, zone int) types.ESSEvent {
	return types.ESSEvent{
		EventName: "GainExperience",
		WorldID:   uint16(world),
		ZoneID:    uint32(zone),

		CharacterID: getRandomCharacterID(),
		LoadoutID:   mkRandomLoadout(),
		TeamID:      mkRandomFaction(),

		ExperienceID: rand.Uint32() % 256,
	}
}

func getRandomCharacterID() string {
	i := rand.Intn(PreloadedCharacterCount)
	return characterStore[i]
}

func mkRandomCharacterID() string {
	return strconv.Itoa(rand.Int())
}

func mkRandomFaction() types.Faction {
	return types.Faction(rand.Intn(4))
}

func mkRandomLoadout() uint16 {
	for {
		i := rand.Intn(46)
		_, ok := translators.LoadoutMap[uint16(i)]
		if ok {
			return uint16(i)
		}
	}
}
