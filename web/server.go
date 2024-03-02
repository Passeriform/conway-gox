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
	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/patterns"
	"github.com/passeriform/conway-gox/web/session"
)

const (
	gameTick = 300 * time.Millisecond
)

var (
	tmpl  map[string]*template.Template   = make(map[string]*template.Template)
	games map[string]*session.GameSession = make(map[string]*session.GameSession)
)

func getTemplatePaths(tmplNames ...string) []string {
	tmplPaths := make([]string, len(tmplNames))
	for idx, name := range tmplNames {
		tmplPaths[idx] = filepath.Join(getServerDir(), "templates", fmt.Sprintf("%v.tmpl", name))
	}
	return tmplPaths
}

func generateTemplate(tmplName string, fnMap template.FuncMap) (*template.Template, error) {
	if tmpl[tmplName] != nil {
		return tmpl[tmplName], nil
	}

	newTmpl := template.New(tmplName)
	if fnMap != nil {
		newTmpl.Funcs(fnMap)
	}
	switch tmplName {
	case "landing":
		return newTmpl.ParseFiles(getTemplatePaths("shell", "landing", "heading")...)
	case "game":
		return newTmpl.ParseFiles(getTemplatePaths("shell", "game", "heading")...)
	case "gameSwap":
		return newTmpl.ParseFiles(getTemplatePaths("game", "heading")...)
	default:
		return nil, fmt.Errorf("unknown template requested: %v", tmplName)
	}
}

func getServerDir() string {
	e, ok := os.LookupEnv("ENVIRONMENT")
	if ok && e != "DEVELOPMENT" {
		return "./"
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "Could not fetch runtime caller to get server directory.")
		os.Exit(1)
	}
	return filepath.Dir(filename)
}

func spawnGame(pattern string) (session.GameSession, error) {
	cellMap := cell_map.New()
	cells, err := patterns.GetPrimitive(pattern, 0, 0)
	if err != nil {
		return session.GameSession{}, fmt.Errorf("unable to fetch the primitive pattern: %v", err)
	}
	cellMap.AddCells(cells)
	return session.NewGameSession(cellMap, session.GameSessionConfiguration{Tick: gameTick}), nil
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("landing", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v", err)
		os.Exit(1)
	}
	if err := tmpl.ExecuteTemplate(w, "shell", patterns.GetAvailablePatterns()); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v", err)
		os.Exit(1)
	}
}

func gameSwapHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("id")

	_, found := games[gameId]

	if !found {
		newGameHandler(w, r)
		return
	}

	gameSession := games[gameId]

	tmpl, err := generateTemplate("gameSwap", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v", err)
		os.Exit(1)
	}

	if err := tmpl.ExecuteTemplate(w, "page", gameSession); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v", err)
		os.Exit(1)
	}
}

// TODO: Add a spinner loader from /connect/ call until canvas first blit is done
// TODO: Give light/dark mode
// TODO: Add ability to draw on canvas
// TODO: Add save and load states
// TODO: Add rewind functionality
// TODO: Fix socket drop and add reconnection logic
// TODO: Add UI messages for server connection/socket state/iteration counter
// TODO: Fix Play/Pause UI button text change

func gameViewHandler(w http.ResponseWriter, r *http.Request) {
	if h := r.Header["Hx-Request"]; h != nil {
		gameSwapHandler(w, r)
		return
	}

	gameId := r.PathValue("id")

	_, found := games[gameId]

	if !found {
		newGameHandler(w, r)
		return
	}

	gameSession := games[gameId]

	tmpl, err := generateTemplate("game", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v", err)
		os.Exit(1)
	}

	// TODO: Set initial state for encode JSON from template directly once sessioned games are implemented.
	if err := tmpl.ExecuteTemplate(w, "shell", gameSession); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v", err)
		os.Exit(1)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		pattern = "PentaDecathlon"
	}
	game, err := spawnGame(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to spawn a new game: %v", err)
		http.Error(w, fmt.Sprintf("Could not initialize game: %v", err), http.StatusInternalServerError)
	}
	games[game.Id] = &game
	http.Redirect(w, r, fmt.Sprintf("/game/%v", game.Id), http.StatusFound)
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
	ioSocket, listenerChannel, err := io.NewSocket(w, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize socket IO handler: %v\n", err)
		http.Error(w, fmt.Sprintf("Could not initialize socket: %v", err), http.StatusInternalServerError)
		return
	}

	eventHandler := func(listener <-chan io.SocketMessage) {
		for e := range listener {
			switch e.Action {
			case "loadState":
			case "saveState":
			case "togglePause":
				gameSession.Game.Running = !gameSession.Game.Running
				ioSocket.MessageChannel <- io.SocketMessage{Action: "pauseToggled", Payload: gameSession.Game.Running}
			case "step":
				if gameSession.Game.Running {
					gameSession.Game.Running = false
					ioSocket.MessageChannel <- io.SocketMessage{Action: "pauseToggled", Payload: gameSession.Game.Running}
				}
				gameSession.Game.Step()
			}
		}
	}

	gameSession.ConnectIO(&ioSocket, listenerChannel, eventHandler)

	go ioSocket.ListenEvents()
	go ioSocket.SendMessages()
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
	mux.Handle("GET /", http.HandlerFunc(landingHandler))

	host := "localhost"
	e, ok := os.LookupEnv("ENVIRONMENT")
	if ok && e != "DEVELOPMENT" {
		host = ""
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	fmt.Fprintf(os.Stdout, "Starting server at %v:%v\n", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), mux); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch the server: %v\n", err)
		os.Exit(1)
	}
}
