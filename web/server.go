package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/patterns"
)

var updateRate = time.Duration(300) * time.Millisecond
var currentGame game.Game
var tmpl *template.Template

func getServerDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return filepath.Dir(filename)
}

func gameViewHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl == nil {
		tmpl = template.Must(template.New("index").Funcs(template.FuncMap{
			"Tick": func() int64 {
				return updateRate.Milliseconds()
			},
			"EncodeJson": func(padding int) [][2]int {
				return currentGame.State.EncodeJson(padding)
			},
		}).ParseFiles(
			filepath.Join(getServerDir(), "templates", "index.tmpl"),
		))
	}
	if err := tmpl.ExecuteTemplate(w, "index", currentGame); err != nil {
		panic(err)
	}
}

// TODO: Use socket instead of polling on client
// TODO: Use game channels and push update to the socket
// TODO: Add multiple game spawner based on id

func stateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currentGame.State.EncodeJson(10)); err != nil {
		panic(err)
	}
}

func nextStepHandler(w http.ResponseWriter, r *http.Request) {
	currentGame.State.Step()
	stateHandler(w, r)
}

func main() {
	var cellMap = cell_map.Create()
	var cells = patterns.GetPrimitive("PentaDecathlon", 0, 0)
	cellMap.AddCells(cells)
	currentGame = game.Create(cellMap, time.Tick(updateRate))

	staticFs := http.FileServer(http.Dir(filepath.Join(getServerDir(), "static")))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.Handle("/game/", http.HandlerFunc(gameViewHandler))
	mux.Handle("/state/", http.HandlerFunc(stateHandler))
	mux.Handle("/step/", http.HandlerFunc(nextStepHandler))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
}
