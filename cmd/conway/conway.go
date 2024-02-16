package main

import (
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
)

func generateMap() cell_map.Map {
	cellMap := cell_map.Create()
	cells := patterns.GetPrimitive("PentaDecathlon", 0, 0)
	cellMap.AddCells(cells)
	return cellMap
}

func main() {
	// Game Creation
	game := game.Create(generateMap(), time.Tick(100*time.Millisecond))

	// IO Handler
	bounds := (<-game.State).GetBounds()
	ioHandler := io.Create(bounds.Right-bounds.Left, bounds.Bottom-bounds.Top, 1)
	defer ioHandler.Close()
	events := make(chan tcell.Event, 1)
	defer close(events)

	// Go Routines
	go game.Play()
	go ioHandler.ListenEvents(events)

	// Channel Handler
	for {
		select {
		case s := <-game.State:
			ioHandler.Blit(s)
		case ev := <-events:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Rune() {
				case 'q':
					game.Exiting <- true
				}
			}
		case e := <-game.Exiting:
			if e {
				os.Exit(0)
			}
		}
	}
}
