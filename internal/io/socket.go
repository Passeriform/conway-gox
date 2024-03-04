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
	conn            *websocket.Conn
	MessageChannel  chan SocketMessage
	listenerChannel chan SocketMessage
	done            chan struct{}
	once            sync.Once
}

func NewSocket(w http.ResponseWriter, r *http.Request) (Socket, <-chan SocketMessage, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while initializing socket: %v\n", err)
		return Socket{}, nil, err
	}
	lChan := make(chan SocketMessage)
	return Socket{conn: conn, MessageChannel: make(chan SocketMessage), listenerChannel: lChan, done: make(chan struct{})}, lChan, nil
}

func (s *Socket) beforeRead() {
	s.conn.SetReadLimit(maxMessageSize)
	readDeadlineFunc := func(string) error {
		if err := s.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return fmt.Errorf("unable to set read deadline on socket: %v", err)
		}
		return nil
	}
	if err := readDeadlineFunc(""); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set read deadline on socket: %v\n", err)
		return
	}
	s.conn.SetPongHandler(readDeadlineFunc)
}

func (s *Socket) beforeWrite() {
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
}

func (s *Socket) SendMessages() {
	defer s.Close()
	for {
		select {
		case <-s.done:
			return
		case message, ok := <-s.MessageChannel:
			if !ok {
				return
			}
			s.beforeWrite()
			jsonBytes, err := json.Marshal(message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to marshal message: %v", message)
				return
			}
			if err := s.conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
				fmt.Fprintf(os.Stderr, "failed at writing message to the writer: %v", err)
				return
			}
		}
	}
}

func (s *Socket) Blit(mapChannel <-chan cell_map.Map) {
	defer s.Close()
	for {
		select {
		case <-s.done:
			return
		case cm, ok := <-mapChannel:
			if !ok {
				return
			}
			s.MessageChannel <- SocketMessage{"updateState", cm.EncodeJson(0)}
		}
	}
}

func (s *Socket) ListenEvents() {
	defer s.Close()
	for {
		select {
		case <-s.done:
			return
		default:
			s.beforeRead()
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to read message on socket: %v\n", err)
				return
			}
			if len(message) == 0 {
				fmt.Println("Received a heartbeat message.")
				s.listenerChannel <- SocketMessage{Action: "heartbeat"}
				continue
			}
			var messageObject SocketMessage
			if err := json.Unmarshal(message, &messageObject); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to unmarshal message: %v\n", message)
				return
			}
			s.listenerChannel <- messageObject
		}
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
		s.done <- struct{}{}
		s.listenerChannel <- SocketMessage{Action: "close"}
	})
}
