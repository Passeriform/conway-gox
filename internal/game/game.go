package game

import (
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

type Game struct {
	State   chan cell_map.Map
	Running chan bool
	Exiting chan bool
	tick    <-chan time.Time
}

func Create(m cell_map.Map, t <-chan time.Time) Game {
	state := make(chan cell_map.Map, 1)
	defer close(state)
	running := make(chan bool, 1)
	defer close(running)
	exiting := make(chan bool, 1)
	defer close(exiting)
	return Game{state, running, exiting, t}
}

func (g *Game) Play() {
	for {
		select {
		case <-g.tick:
			if <-g.Running {
				currentState := <-g.State
				currentState.Step()
				g.State <- currentState
			}
		}
	}
}
