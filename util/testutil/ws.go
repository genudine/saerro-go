package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
)

type MockESS struct {
	Server      *httptest.Server
	LastMessage string
}

func (m MockESS) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c, err := websocket.Accept(rw, req, nil)
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte("websocket connection failed"))
		return
	}
	defer c.CloseNow()

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*30)
	defer cancel()

	_, body, err := c.Read(ctx)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte("websocket read failed"))
		return
	}

	m.LastMessage = string(body)
}

func GetMockESS(t *testing.T) MockESS {
	t.Helper()

	m := MockESS{}

	s := httptest.NewServer(m)
	m.Server = s

	t.Cleanup(s.Close)
	return m
}
