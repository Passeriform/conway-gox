package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
)

func main() {
	// Map Creation
	cellMap := cell_map.New()
	cells := patterns.GetPrimitive("Toad", 0, 0)
	cellMap.AddCells(cells)

	// Game Creation
	game := game.New(cellMap, time.Tick(100*time.Millisecond))

	// IO Handler
	ioHandler, err := io.NewTerminal()
	if err != nil {
		fmt.Printf("Could not initialize terminal IO handler: %v", err)
		return
	}
	defer ioHandler.Close()

	// Go Routines
	stateChannel := make(chan cell_map.Map, 1)
	eventChannel := make(chan tcell.Event, 1)
	go game.Play(stateChannel)
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
