package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genudine/saerro-go/cmd/ws/eventhandler"
	"github.com/genudine/saerro-go/cmd/ws/wsmanager"
	"github.com/genudine/saerro-go/types"
)

func main() {
	wsAddr := os.Getenv("WS_ADDR")
	if wsAddr == "" {
		log.Fatalln("WS_ADDR is not set.")
	}

	wsm := wsmanager.NewWebsocketManager(eventhandler.EventHandler{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := wsm.Connect(ctx, wsAddr)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.Buffer{}

	go func() {
		for {
			ctx := context.Background()

			_, data, err := wsm.Conn.Read(ctx)
			if err != nil {
				log.Fatalln("wsm: read failed:", err)
			}

			buf.Write(data)
			buf.WriteByte('\n')
		}
	}()

	go func() {
		time.Sleep(time.Second * 1)
		err = wsm.Subscribe(ctx)
		if err != nil {
			wsm.FailClose()
			log.Fatalln("subscribe failed", err)
		}
		// log.Println("sent subscribe")
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-exitSignal:
		log.Println("got interrupt, exiting...")
	case <-wsm.Closed:
		log.Println("websocket closed, bailing...")
	case <-time.After(time.Second * 30):
		buf.WriteTo(os.Stdout)
	}
}

type handler struct{}

func (h handler) HandleEvent(_ context.Context, _ types.ESSEvent) {}
