package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
)

func main() {
	// Map Creation
	cellMap := cell_map.Create()
	cells := patterns.GetPrimitive("Toad", 0, 0)
	cellMap.AddCells(cells)

	// Game Creation
	game := game.Create(cellMap, time.Tick(100*time.Millisecond))
	defer game.Close()

	// IO Handler
	ioHandler := io.Create()
	defer ioHandler.Close()
	events := make(chan tcell.Event, 1)
	defer close(events)

	// Go Routines
	go game.Play()
	go ioHandler.ListenEvents(events)

	// Channel Handler
	for {
		select {
		case s := <-game.StateChange:
			ioHandler.Blit(s)
		case ev := <-events:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Rune() {
				case 'q':
					game.Close()
				case 'p':
					game.Running = !game.Running
				}
			}
		case <-game.Exit:
			break
		}
	}
}
