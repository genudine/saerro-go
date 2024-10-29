package eventhandler

import (
	"context"
	"testing"
	"time"

	"github.com/genudine/saerro-go/cmd/ws/ingest"
	"github.com/genudine/saerro-go/store/storemock"
	"github.com/genudine/saerro-go/types"
)

func getEventHandlerTestShim(t *testing.T) (EventHandler, context.Context, *storemock.MockPlayerStore) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	t.Cleanup(cancel)

	ps := new(storemock.MockPlayerStore)

	return EventHandler{
		Ingest: &ingest.Ingest{
			PlayerStore: ps,
		},
	}, ctx, ps
}

func TestHandleDeath(t *testing.T) {
	eh, ctx, ps := getEventHandlerTestShim(t)

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

	p1 := types.PopEventFromESSEvent(event, false).ToPlayer()
	p2 := types.PopEventFromESSEvent(event, true).ToPlayer()

	ps.On("Insert", ctx, p1).Return(nil)
	ps.On("Insert", ctx, p2).Return(nil)

	eh.HandleDeath(ctx, event)
}

func TestHandleExperience(t *testing.T) {
	eh, ctx, ps := getEventHandlerTestShim(t)

	event := types.ESSEvent{
		EventName: "GainExperience",
		WorldID:   17,
		ZoneID:    2,

		CharacterID: "LyytisDoll",
		LoadoutID:   3,
		TeamID:      types.NC,

		ExperienceID: 674,
	}

	p := types.PopEventFromESSEvent(event, false).ToPlayer()
	ps.On("Insert", ctx, p).Return(nil)

	eh.HandleExperience(ctx, event)
}

func TestHandleAnalytics(t *testing.T) {
	eh, ctx, _ := getEventHandlerTestShim(t)
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
	eh, ctx, ps := getEventHandlerTestShim(t)

	events := []types.ESSEvent{
		{
			EventName: "Death",
			WorldID:   17,
			ZoneID:    2,

			CharacterID:        "LyytisDoll",
			CharacterLoadoutID: 3,
			TeamID:             types.NC,

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
		{
			EventName: "",
		},
	}

	p1 := types.PopEventFromESSEvent(events[0], false).ToPlayer()
	ps.On("Insert", ctx, p1).Return(nil).Once()
	p2 := types.PopEventFromESSEvent(events[0], true).ToPlayer()
	ps.On("Insert", ctx, p2).Return(nil).Once()
	p3 := types.PopEventFromESSEvent(events[1], false).ToPlayer()
	ps.On("Insert", ctx, p3).Return(nil).Once()

	for _, event := range events {
		eh.HandleEvent(ctx, event)
	}
}
