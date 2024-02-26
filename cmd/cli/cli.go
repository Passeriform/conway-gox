package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
)

// TODO: Cli is broken. Fix required

func main() {
	// Map Creation
	cellMap := cell_map.New()
	cells, err := patterns.GetPrimitive("Toad", 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while fetching primitive: %v", err)
	}
	cellMap.AddCells(cells)

	// Game Creation
	game, stateChannel := game.New(cellMap, time.Tick(100*time.Millisecond))

	// IO Handler
	ioHandler, err := io.NewTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize terminal IO handler: %v\n", err)
		return
	}
	defer ioHandler.Close()

	// Go Routines
	eventChannel := make(chan tcell.Event)
	go game.Play()
	go ioHandler.Blit(stateChannel)
	go ioHandler.ListenEvents(eventChannel)
	for e := range eventChannel {
		switch e := e.(type) {
		case *tcell.EventKey:
			switch e.Rune() {
			case 'p':
				game.Running = !game.Running
			}
		}
	}
}
