package session

import (
	"math/rand"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/web/multiplexer"
)

const (
	letterBytes         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	gameSessionIdLength = 5
)

type GameSessionConfiguration struct {
	Tick time.Duration
}

type GameSession struct {
	Id           string
	Game         game.Game
	stateMux     *multiplexer.Multiplexer[cell_map.Map]
	eventChannel chan io.SocketMessage
}

func generateGameId() string {
	b := make([]byte, gameSessionIdLength)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func NewGameSession(initialMap cell_map.Map, eventHandler func(eventChannel <-chan io.SocketMessage, state *game.Game), config GameSessionConfiguration) GameSession {
	id := generateGameId()
	eventChannel := make(chan io.SocketMessage)
	newGame, stateChannel := game.New(initialMap, time.Tick(config.Tick))
	stateMux := multiplexer.Of(stateChannel)
	go stateMux.Forward()
	go newGame.Play()
	go eventHandler(eventChannel, &newGame)
	return GameSession{id, newGame, &stateMux, eventChannel}
}

func (g *GameSession) ConnectIO(ioHandler *io.Socket) {
	go ioHandler.Blit(g.stateMux.ProvisionSink())
	go ioHandler.ListenEvents(g.eventChannel)
}
