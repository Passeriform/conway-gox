package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
)

var (
	tmpl        *template.Template
	currentGame game.Game
)

func getServerDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return filepath.Dir(filename)
}

func gameViewHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl == nil {
		// TODO: Set initial state for encode JSON from template directly once sessioned games are implemented.
		tmpl = template.Must(template.New("index").ParseFiles(
			filepath.Join(getServerDir(), "templates", "index.tmpl"),
		))
	}
	if err := tmpl.ExecuteTemplate(w, "index", currentGame); err != nil {
		panic(err)
	}
}

func spawnGameHandler(w http.ResponseWriter, r *http.Request) {
	// Map Creation
	cellMap := cell_map.New()
	cells := patterns.GetPrimitive("PentaDecathlon", 0, 0)
	cellMap.AddCells(cells)

	// Game Creation
	currentGame = game.New(cellMap, time.Tick(300*time.Millisecond))

	// IO Handler
	ioHandler, err := io.NewSocket(w, r)
	if err != nil {
		fmt.Printf("Could not initialize socket IO handler: %v", err)
		return
	}
	defer ioHandler.Close()

	// Go Routines
	stateChannel := make(chan cell_map.Map, 1)
	eventChannel := make(chan io.SocketMessage, 1)
	go currentGame.Play(stateChannel)
	go ioHandler.Blit(stateChannel)
	go ioHandler.ListenEvents(eventChannel)
	for e := range eventChannel {
		switch e.Action {
		case "loadState":
		case "saveState":
		case "togglePause":
			currentGame.Running = !currentGame.Running
		case "step":
			currentGame.Running = false
			currentGame.State.Step()
			stateChannel <- *currentGame.State
		}
	}
}

// TODO: Add multiple game spawner based on id

func main() {
	defer func() {
		// TODO: Implement if required
	}()

	staticFs := http.FileServer(http.Dir(filepath.Join(getServerDir(), "static")))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.Handle("/game/", http.HandlerFunc(gameViewHandler))
	mux.Handle("/connect/", http.HandlerFunc(spawnGameHandler))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
}
