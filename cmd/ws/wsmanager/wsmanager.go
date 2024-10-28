package wsmanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/genudine/saerro-go/cmd/ws/eventhandler"
	"github.com/genudine/saerro-go/types"
)

type WebsocketManager struct {
	Conn         *websocket.Conn
	EventHandler eventhandler.EventHandler
	Closed       chan bool
}

func NewWebsocketManager(eh eventhandler.EventHandler) WebsocketManager {
	return WebsocketManager{
		EventHandler: eh,
		Closed:       make(chan bool, 1),
	}
}

func (wsm *WebsocketManager) Connect(ctx context.Context, addr string) (err error) {
	wsm.Conn, _, err = websocket.Dial(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("wsm: connect failed: %w", err)
	}

	log.Println("wsm: connected to", addr)

	return
}

type ESSSubscription struct {
	Action                         string   `json:"action,omitempty"`
	Worlds                         []string `json:"worlds,omitempty"`
	EventNames                     []string `json:"eventNames,omitempty"`
	Characters                     []string `json:"characters,omitempty"`
	Service                        string   `json:"service,omitempty"`
	LogicalAndCharactersWithWorlds bool     `json:"logicalAndCharactersWithWorlds,omitempty"`
}

func (wsm *WebsocketManager) Subscribe(ctx context.Context) error {
	sub := ESSSubscription{
		Action:                         "subscribe",
		Service:                        "event",
		Worlds:                         []string{"all"},
		EventNames:                     getEventNames(),
		Characters:                     []string{"all"},
		LogicalAndCharactersWithWorlds: true,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(sub)
	if err != nil {
		return fmt.Errorf("wsm: subscribe: json encode failed: %w", err)
	}

	log.Printf("wsm: subscribe message: %s", buf.String())

	err = wsm.Conn.Write(ctx, websocket.MessageText, buf.Bytes())
	if err != nil {
		return fmt.Errorf("wsm: subscribe: ws write failed: %w", err)
	}

	return nil
}

func (wsm *WebsocketManager) Start() {
	go wsm.startWatchdog()

	for {
		ctx := context.Background()

		var event types.ESSData

		_, data, err := wsm.Conn.Read(ctx)
		if err != nil {
			log.Fatalln("wsm: read failed:", err)
		}

		// log.Printf("raw event: %s", string(data))

		err = json.Unmarshal(data, &event)
		if err != nil {
			log.Println("wsm: json unmarshal failed:", err)
			log.Println("wsm: json unmarshal failed (payload)", string(data))
		}

		go wsm.EventHandler.HandleEvent(ctx, event.Payload)
	}
}

func (wsm *WebsocketManager) startWatchdog() {
	for {
		time.Sleep(time.Second * 30)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err := wsm.Conn.Ping(ctx)
		if err != nil {
			log.Println("wsm: watchdog failed")
			wsm.Closed <- true
		}

		cancel()
	}
}

func (wsm *WebsocketManager) Close() {
	wsm.Conn.Close(websocket.StatusNormalClosure, "")
	wsm.Closed <- true
}

func (wsm *WebsocketManager) FailClose() {
	wsm.Conn.Close(websocket.StatusAbnormalClosure, "")
	wsm.Closed <- true
}
