package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/genudine/saerro-go/store"
	"github.com/genudine/saerro-go/util"
)

func main() {
	db, err := util.GetDBConnection(os.Getenv("DB_ADDR"))
	if err != nil {
		log.Fatalln(err)
	}

	playerStore := store.NewPlayerStore(db)
	vehicleStore := store.NewVehicleStore(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	run(ctx, playerStore, vehicleStore)
}

func run(ctx context.Context, ps store.IPlayerStore, vs store.IVehicleStore) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i, err := ps.Prune(ctx)
		if err != nil {
			log.Println("pruner: playerStore.Prune failed")
		}

		log.Printf("pruner: deleted %d players", i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i, err := vs.Prune(ctx)
		if err != nil {
			log.Println("pruner: vehicleStore.Prune failed")
		}

		log.Printf("pruner: deleted %d vehicles", i)
	}()

	wg.Wait()
}
