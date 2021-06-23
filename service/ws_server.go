package main

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WsReadWriteCloser struct {
	c *websocket.Conn
}

func (rwc *WsReadWriteCloser) Read(b []byte) (int, error) {
	_, out, err := rwc.c.ReadMessage()
	copy(b, out)
	return len(out), err
}

func (rwc *WsReadWriteCloser) Write(b []byte) (int, error) {
	err := rwc.c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (rwc *WsReadWriteCloser) Close() error {
	return rwc.c.Close()
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	rwc := &WsReadWriteCloser{conn}

	// TODO reuse session ID?
	s := NewSession(context.TODO(), GenerateUUID(), rwc, rwc)
	SessionManager.AddSession(s)

	s.Loop()
}
