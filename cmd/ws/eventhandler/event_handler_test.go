package eventhandler

import (
	"context"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"

	"github.com/genudine/saerro-go/cmd/ws/ingest"
	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/translators"
	"github.com/genudine/saerro-go/types"
	"github.com/genudine/saerro-go/util/testutil"
)

func getEventHandlerTestShim(t *testing.T) (EventHandler, context.Context) {
	t.Helper()

	db := testutil.GetTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	t.Cleanup(cancel)

	return EventHandler{
		Ingest: &ingest.Ingest{
			PlayerStore: store.NewPlayerStore(db),
		},
	}, ctx
}

func TestHandleDeath(t *testing.T) {
	eh, ctx := getEventHandlerTestShim(t)

	event := types.ESSEvent{
		EventName: "Death",
		WorldID:   17,
		ZoneID:    2,

		CharacterID: "LyytisDoll",
		LoadoutID:   3,
		TeamID:      types.NC,

		AttackerCharacterID: "Lyyti",
		AttackerLoadoutID:   3,
		AttackerTeamID:      types.TR,
	}

	eh.HandleDeath(ctx, event)

	player1, err := eh.Ingest.PlayerStore.GetOne(ctx, event.CharacterID)
	assert.NoError(t, err, "player1 fetch failed")
	assert.Equal(t, event.CharacterID, player1.CharacterID)
	assert.Equal(t, string(translators.ClassFromLoadout(event.LoadoutID)), player1.ClassName)

	player2, err := eh.Ingest.PlayerStore.GetOne(ctx, event.AttackerCharacterID)
	assert.NoError(t, err, "player2 fetch failed")
	assert.Equal(t, event.AttackerCharacterID, player2.CharacterID)
	assert.Equal(t, string(translators.ClassFromLoadout(event.AttackerLoadoutID)), player2.ClassName)
}

func TestHandleExperience(t *testing.T) {
	eh, ctx := getEventHandlerTestShim(t)

	event := types.ESSEvent{
		EventName: "GainExperience",
		WorldID:   17,
		ZoneID:    2,

		CharacterID: "LyytisDoll",
		LoadoutID:   3,
		TeamID:      types.NC,

		ExperienceID: 674,
	}

	eh.HandleExperience(ctx, event)
	player, err := eh.Ingest.PlayerStore.GetOne(ctx, event.CharacterID)
	assert.NoError(t, err, "player fetch check failed")
	assert.Equal(t, event.CharacterID, player.CharacterID)
	assert.Equal(t, string(translators.ClassFromLoadout(event.LoadoutID)), player.ClassName)
}

func TestHandleAnalytics(t *testing.T) {
	eh, ctx := getEventHandlerTestShim(t)
	event := types.ESSEvent{
		EventName: "GainExperience",
		WorldID:   17,
		ZoneID:    2,

		CharacterID: "LyytisDoll",
		LoadoutID:   3,
		TeamID:      types.NC,

		ExperienceID: 674,
	}

	eh.HandleAnalytics(ctx, event)
}

func TestHandleEvent(t *testing.T) {
	eh, ctx := getEventHandlerTestShim(t)

	events := []types.ESSEvent{
		{
			EventName: "Death",
			WorldID:   17,
			ZoneID:    2,

			CharacterID: "LyytisDoll",
			LoadoutID:   3,
			TeamID:      types.NC,

			AttackerCharacterID: "Lyyti",
			AttackerLoadoutID:   3,
			AttackerTeamID:      types.TR,
		},
		{
			EventName: "GainExperience",
			WorldID:   17,
			ZoneID:    2,

			CharacterID: "DollNC",
			LoadoutID:   3,
			TeamID:      types.NC,

			ExperienceID: 201,
		},
	}

	for _, event := range events {
		eh.HandleEvent(ctx, event)
	}

	checkPlayers := []string{"LyytisDoll", "Lyyti", "DollNC"}
	for _, id := range checkPlayers {
		// eventual consistency <333
		err := retry.Do(func() error {
			_, err := eh.Ingest.PlayerStore.GetOne(ctx, id)
			return err
		})
		assert.NoError(t, err)
	}
}
