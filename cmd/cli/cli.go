package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/loader"
)

// TODO: Cli is broken. Fix required

func main() {
	// Map Creation
	primitive, err := loader.LoadFromPrimitive("toad", 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load state from primitive pattern: %v\n", err)
		return
	}

	// Game Creation
	game, stateChannel := game.New(primitive, time.Tick(100*time.Millisecond))

	// IO Handler
	ioHandler, listenerChannel, err := io.NewTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize terminal IO handler: %v\n", err)
		return
	}
	defer ioHandler.Close()

	// Go Routines
	go game.Play()
	go ioHandler.SendMessages()
	go ioHandler.Blit(stateChannel)
	go ioHandler.ListenEvents()
	for e := range listenerChannel {
		switch e := e.(type) {
		case *tcell.EventKey:
			switch e.Rune() {
			case 'p':
				game.Running = !game.Running
			}
		}
	}
}
