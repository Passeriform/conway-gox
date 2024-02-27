package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/game"
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
	"github.com/passeriform/conway-gox/web/session"
)

const (
	gameTick = 300 * time.Millisecond
	host     = "localhost:8080"
)

var (
	tmpl  *template.Template
	games map[string]*session.GameSession = make(map[string]*session.GameSession)
)

func getServerDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "Could not fetch runtime caller to get server directory.")
		os.Exit(1)
	}
	return filepath.Dir(filename)
}

func spawnGame() (session.GameSession, error) {
	cellMap := cell_map.New()
	cells, err := patterns.GetPrimitive("PentaDecathlon", 0, 0)
	if err != nil {
		return session.GameSession{}, fmt.Errorf("unable to fetch the primitive pattern: %v", err)
	}
	cellMap.AddCells(cells)
	eventHandler := func(input <-chan io.SocketMessage, currentGame *game.Game) {
		for e := range input {
			switch e.Action {
			case "loadState":
			case "saveState":
			case "togglePause":
				currentGame.Running = !currentGame.Running
			case "step":
				currentGame.Running = false
				currentGame.Step()
			}
		}
	}
	return session.NewGameSession(cellMap, eventHandler, session.GameSessionConfiguration{Tick: gameTick}), nil
}

func gameViewHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("id")

	_, found := games[gameId]

	if !found {
		newGameHandler(w, r)
		return
	}

	gameSession := games[gameId]

	if tmpl == nil {
		// TODO: Set initial state for encode JSON from template directly once sessioned games are implemented.
		tmpl = template.Must(template.New("index").ParseFiles(
			filepath.Join(getServerDir(), "templates", "index.tmpl"),
		))
	}
	if err := tmpl.ExecuteTemplate(w, "index", gameSession); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v", err)
		os.Exit(1)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	game, err := spawnGame()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to spawn a new game: %v", err)
		http.Error(w, fmt.Sprintf("Could not initialize game: %v", err), http.StatusInternalServerError)
	}
	games[game.Id] = &game
	http.Redirect(w, r, "/game/"+game.Id, http.StatusSeeOther)
}

func connectClientHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch Game
	gameId := r.PathValue("id")
	gameSession, found := games[gameId]
	if !found {
		fmt.Fprintf(os.Stderr, "Requested a game that doesn't exist: %v\n", gameId)
		http.Error(w, fmt.Sprintf("Requested a game that doesn't exist: %v", gameId), http.StatusNotFound)
		return
	}

	// IO Handler
	ioSocket, err := io.NewSocket(w, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize socket IO handler: %v\n", err)
		http.Error(w, fmt.Sprintf("Could not initialize socket: %v", err), http.StatusInternalServerError)
		return
	}

	gameSession.ConnectIO(&ioSocket)
}

func main() {
	defer func() {
		// TODO: Check for leaked resources, channels, websockets, gamesession and multiplexer
		// TODO: Implement if required
	}()

	staticFs := http.FileServer(http.Dir(filepath.Join(getServerDir(), "static")))

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFs))
	mux.Handle("GET /connect/{id}", http.HandlerFunc(connectClientHandler))
	mux.Handle("GET /game/{id}", http.HandlerFunc(gameViewHandler))
	mux.Handle("GET /game/", http.HandlerFunc(newGameHandler))
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/game/", http.StatusSeeOther)
	}))

	fmt.Fprintf(os.Stdout, "Starting server at %v\n", host)
	if err := http.ListenAndServe(host, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch the server: %v\n", err)
		os.Exit(1)
	}
}
