package game

import (
	"fmt"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/loader"
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
	g.state.Next()
	g.stateChannel <- *g.state
}

func (g *Game) SaveState(saveFp string, padding int) error {
	if err := loader.SaveToFile(*g.state, saveFp, padding); err != nil {
		return fmt.Errorf("unable to save game state: %v", err)
	}
	return nil
}

func (g *Game) LoadState(saveFp string, padding int) error {
	state, err := loader.LoadFromFile(saveFp, padding)
	if err != nil {
		return fmt.Errorf("unable to load game state: %v", err)
	}
	g.state = &state
	g.stateChannel <- *g.state
	return nil
}

// TODO: Add close method with sync.once, ensure stateChannel is not listened to once closed
