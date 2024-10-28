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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		playerStore := store.NewPlayerStore(db)
		i, err := playerStore.Prune(ctx)
		if err != nil {
			log.Println("pruner: playerStore.Prune failed")
		}

		log.Printf("pruner: deleted %d players", i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		vehicleStore := store.NewVehicleStore(db)
		i, err := vehicleStore.Prune(ctx)
		if err != nil {
			log.Println("pruner: vehicleStore.Prune failed")
		}

		log.Printf("pruner: deleted %d vehicles", i)
	}()

	wg.Wait()
}
