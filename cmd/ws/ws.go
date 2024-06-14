package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/genudine/saerro-go/types"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	wsAddr := os.Getenv("WS_ADDR")
	if wsAddr == "" {
		log.Fatalln("WS_ADDR is not set.")
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalln("database connection failed", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	wsConn, _, err := websocket.Dial(ctx, wsAddr, nil)
	if err != nil {
		log.Fatalln("Connection to ESS failed.", err)
	}
	defer wsConn.Close(websocket.StatusInternalError, "internal error. bye")

	err = wsjson.Write(ctx, wsConn, map[string]interface{}{
		"action":     "subscribe",
		"worlds":     "all",
		"eventNames": getEventNames(),
		"characters": []string{"all"},
		"service":    "event",

		"logicalAndCharactersWithWorlds": true,
	})
	if err != nil {
		log.Fatalln("subscription write failed", err)
	}

	log.Println("subscribe done")

	eventHandler := EventHandler{
		Ingest: &Ingest{
			DB: db,
		},
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		var event types.ESSData
		err := wsjson.Read(ctx, wsConn, &event)
		if err != nil {
			log.Println("wsjson read failed", err)
			cancel()
			continue
		}

		go eventHandler.HandleEvent(ctx, event.Payload)
	}

	wsConn.Close(websocket.StatusNormalClosure, "")
}
