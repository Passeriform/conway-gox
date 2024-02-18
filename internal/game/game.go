package game

import (
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

type Game struct {
	State       *cell_map.Map
	Running     bool
	StateChange chan cell_map.Map
	Exit        chan struct{}
	tick        <-chan time.Time
}

func Create(m cell_map.Map, t <-chan time.Time) Game {
	stateChange := make(chan cell_map.Map)
	exit := make(chan struct{})
	return Game{&m, false, stateChange, exit, t}
}

func (g *Game) Play() {
	g.Running = true
	for {
		select {
		case <-g.tick:
			if g.Running {
				g.State.Step()
				g.StateChange <- *g.State
			}
		}
	}
}

func (g *Game) Close() {
	close(g.StateChange)
	close(g.Exit)
}
