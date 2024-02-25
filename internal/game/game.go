package game

import (
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

type Game struct {
	State   *cell_map.Map
	Running bool
	Tick    <-chan time.Time
}

func New(m cell_map.Map, t <-chan time.Time) Game {
	return Game{&m, false, t}
}

func (g *Game) Play(stateChannel chan<- cell_map.Map) {
	g.Running = true
	for {
		select {
		case <-g.Tick:
			if g.Running {
				g.State.Step()
				stateChannel <- *g.State
			}
		}
	}
}
