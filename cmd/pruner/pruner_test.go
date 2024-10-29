package main

import (
	"context"
	"testing"
	"time"

	"github.com/genudine/saerro-go/store/storemock"
)

func TestRun(t *testing.T) {
	ps := new(storemock.MockPlayerStore)
	vs := new(storemock.MockVehicleStore)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ps.On("Prune", ctx).Return(30, nil)
	vs.On("Prune", ctx).Return(30, nil)

	run(ctx, ps, vs)
}
