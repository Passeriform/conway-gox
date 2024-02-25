package io

// TODO: Break IO into smaller sub-packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/passeriform/conway-gox/internal/cell_map"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Socket struct {
	conn *websocket.Conn
	once sync.Once
}

func NewSocket(w http.ResponseWriter, r *http.Request) (Socket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return Socket{}, err
	}
	return Socket{conn: conn}, nil
}

func (s *Socket) Blit(mapChannel <-chan cell_map.Map) {
	defer func() {
		s.Close()
	}()

	for cm := range mapChannel {
		s.conn.SetWriteDeadline(time.Now().Add(writeWait))
		s.conn.SetPingHandler(func(string) error { s.conn.SetWriteDeadline(time.Now().Add(writeWait)); return nil })
		w, err := s.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		message := cm.EncodeJson(10)
		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("Unable to marshal message %v", message)
		}
		w.Write(jsonBytes)
	}
}

// TODO: Remove all panic calls

func (s *Socket) ListenEvents(eventChannel chan<- SocketMessage) {
	defer func() {
		s.Close()
	}()

	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			return
		}
		var messageObject SocketMessage
		if err := json.Unmarshal(message, &messageObject); err != nil {
			fmt.Printf("Unable to unmarshal message %v", message)
		}
		eventChannel <- messageObject
	}
}

func (s *Socket) Close() {
	s.once.Do(func() {
		s.conn.WriteMessage(websocket.CloseMessage, []byte{})
		s.conn.Close()
	})
}
