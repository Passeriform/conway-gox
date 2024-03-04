package session

import (
	"sync"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/utility"
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
	Id            string
	Game          *game.Game
	stateMux      *multiplexer.Multiplexer[cell_map.Map]
	eventMux      *multiplexer.Multiplexer[io.SocketMessage]
	outputDoneMap map[*io.Socket]string
	ioRefCount    int
	once          sync.Once
}

func NewGameSession(initialMap cell_map.Map, config GameSessionConfiguration) GameSession {
	id := utility.NewRandomString(gameSessionIdLength)
	newGame, sc := game.New(initialMap, time.Tick(config.Tick))
	stateMux := multiplexer.Of(sc)
	go stateMux.Forward()
	go newGame.Play()
	return GameSession{Id: id, Game: &newGame, stateMux: &stateMux, outputDoneMap: make(map[*io.Socket]string)}
}

func (g *GameSession) ConnectIO(ioSocket *io.Socket, listenerChannel <-chan io.SocketMessage, eventHandler func(<-chan io.SocketMessage)) {
	g.ioRefCount = g.ioRefCount + 1
	eventMux := multiplexer.Of(listenerChannel)
	go eventMux.Forward()
	g.eventMux = &eventMux
	stateChanId, stateChan := g.stateMux.ProvisionSink()
	_, eventChan := g.eventMux.ProvisionSink()
	go ioSocket.Blit(stateChan)
	go eventHandler(eventChan)
	g.outputDoneMap[ioSocket] = stateChanId
}

func (g *GameSession) SignalClose(ioSocket *io.Socket) {
	g.stateMux.CloseSink(g.outputDoneMap[ioSocket])
	g.ioRefCount = g.ioRefCount - 1
	if g.ioRefCount != 0 {
		return
	}
	g.once.Do(func() {
		g.Game.Close()
		g.stateMux.Close()
		g.eventMux.Close()
	})
}
