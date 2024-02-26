package game

import (
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

type Game struct {
	state        *cell_map.Map
	Running      bool
	Tick         <-chan time.Time
	stateChannel chan<- cell_map.Map
}

func New(m cell_map.Map, t <-chan time.Time) (Game, chan cell_map.Map) {
	sc := make(chan cell_map.Map)
	return Game{&m, false, t, sc}, sc
}

func (g *Game) Play() {
	g.Running = true
	for range g.Tick {
		if !g.Running {
			continue
		}
		g.Step()
	}
}

func (g *Game) Step() {
	g.state.Step()
	g.stateChannel <- *g.state
}

// TODO: Add close method with sync.once, ensure stateChannel is not listened to once closed
