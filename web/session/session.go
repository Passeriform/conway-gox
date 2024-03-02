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
	Id       string
	Game     *game.Game
	stateMux *multiplexer.Multiplexer[cell_map.Map]
	eventMux *multiplexer.Multiplexer[io.SocketMessage]
}

func generateGameId() string {
	b := make([]byte, gameSessionIdLength)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func NewGameSession(initialMap cell_map.Map, config GameSessionConfiguration) GameSession {
	id := generateGameId()
	newGame, sc := game.New(initialMap, time.Tick(config.Tick))
	stateMux := multiplexer.Of(sc)
	go stateMux.Forward()
	go newGame.Play()
	return GameSession{Id: id, Game: &newGame, stateMux: &stateMux}
}

func (g *GameSession) ConnectIO(ioSocket *io.Socket, listenerChannel <-chan io.SocketMessage, eventHandler func(<-chan io.SocketMessage)) {
	eventMux := multiplexer.Of(listenerChannel)
	go eventMux.Forward()
	g.eventMux = &eventMux
	go ioSocket.Blit(g.stateMux.ProvisionSink())
	go eventHandler(g.eventMux.ProvisionSink())
}
