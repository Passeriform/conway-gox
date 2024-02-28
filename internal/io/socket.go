package io

// TODO: Break IO into smaller sub-packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type SocketMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload,omitempty"`
}

type Socket struct {
	conn *websocket.Conn
	once sync.Once
}

func NewSocket(w http.ResponseWriter, r *http.Request) (Socket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while initializing socket: %v\n", err)
		return Socket{}, err
	}
	return Socket{conn: conn}, nil
}

func (s *Socket) Blit(mapChannel <-chan cell_map.Map) {
	defer func() {
		s.Close()
	}()

	for cm := range mapChannel {
		writeDeadlineFunc := func(string) error {
			if err := s.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to set write deadline on socket: %v\n", err)
				return err
			}
			return nil
		}
		if err := writeDeadlineFunc(""); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to set first write deadline: %v\n", err)
			return
		}
		s.conn.SetPingHandler(writeDeadlineFunc)
		w, err := s.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch connection writer: %v\n", err)
			return
		}
		message := cm.EncodeJson(0)
		jsonBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal message: %v\n", message)
			return
		}
		if _, err := w.Write(jsonBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Failed at writing message to the writer: %v\n", err)
			return
		}
	}
}

func (s *Socket) ListenEvents(eventChannel chan<- SocketMessage) {
	defer func() {
		s.Close()
	}()

	s.conn.SetReadLimit(maxMessageSize)
	readDeadlineFunc := func(string) error {
		if err := s.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return fmt.Errorf("unable to set read deadline on socket: %v", err)
		}
		return nil
	}
	if err := readDeadlineFunc(""); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set first read deadline: %v\n", err)
		return
	}
	s.conn.SetPongHandler(readDeadlineFunc)

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read message on socket: %v\n", err)
			return
		}
		var messageObject SocketMessage
		if err := json.Unmarshal(message, &messageObject); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to unmarshal message: %v\n", message)
			return
		}
		eventChannel <- messageObject
	}
}

func (s *Socket) Close() {
	s.once.Do(func() {
		if err := s.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write close message to the connection in deferred close method: %v", err)
		}
		if err := s.conn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to close the socket connection in deferred close method: %v", err)
		}
	})
}
