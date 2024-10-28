package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genudine/saerro-go/cmd/ws/eventhandler"
	"github.com/genudine/saerro-go/cmd/ws/wsmanager"
	"github.com/genudine/saerro-go/util"
)

func main() {
	wsAddr := os.Getenv("WS_ADDR")
	if wsAddr == "" {
		log.Fatalln("WS_ADDR is not set.")
	}

	db, err := util.GetDBConnection(os.Getenv("DB_ADDR"))
	if err != nil {
		log.Fatalln(err)
	}

	eventHandler := eventhandler.NewEventHandler(db)
	wsm := wsmanager.NewWebsocketManager(eventHandler)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = wsm.Connect(ctx, wsAddr)
	if err != nil {
		log.Fatalln(err)
	}

	go wsm.Start()

	go func() {
		time.Sleep(time.Second * 1)
		err = wsm.Subscribe(ctx)
		if err != nil {
			wsm.FailClose()
			log.Fatalln("subscribe failed", err)
		}
		log.Println("sent subscribe")
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-exitSignal:
		log.Println("got interrupt, exiting...")
	case <-wsm.Closed:
		log.Println("websocket closed, bailing...")
	}
}
