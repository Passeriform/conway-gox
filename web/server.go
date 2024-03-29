package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/passeriform/conway-gox/internal/io"
	"github.com/passeriform/conway-gox/internal/loader"
	"github.com/passeriform/conway-gox/web/session"
)

const (
	gameTick = 300 * time.Millisecond
	// TODO: Fill on server load
	// TODO: Delete directory on server close
)

var (
	savePath = filepath.Join(os.Getenv("GOPATH"), "saves")
)

var (
	tmpl  map[string]*template.Template   = make(map[string]*template.Template)
	games map[string]*session.GameSession = make(map[string]*session.GameSession)
)

func getTemplatePaths(tmplNames ...string) []string {
	tmplPaths := make([]string, len(tmplNames))
	for idx, name := range tmplNames {
		tmplPaths[idx] = filepath.Join(os.Getenv("GOPATH"), "web", "templates", fmt.Sprintf("%v.tmpl", name))
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

func spawnGame(pattern string) (session.GameSession, error) {
	primitive, err := loader.LoadFromPrimitive(pattern, 0)
	if err != nil {
		return session.GameSession{}, fmt.Errorf("unable to load primitive pattern: %v", err)
	}
	// TODO: Allow usage of save and load states here using game ids
	return session.NewGameSession(primitive, session.GameSessionConfiguration{Tick: gameTick}), nil
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("landing", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v\n", err)
		os.Exit(1)
	}
	pm, err := loader.ScanPrimitivesByType()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while scanning for primitives: %v\n", err)
	}
	if err := tmpl.ExecuteTemplate(w, "shell", pm); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v\n", err)
		os.Exit(1)
	}

	if err := tmpl.ExecuteTemplate(w, "page", gameSession); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "An error occurred while generating template: %v\n", err)
		os.Exit(1)
	}

	// TODO: Set initial state for encode JSON from template directly once sessioned games are implemented.
	if err := tmpl.ExecuteTemplate(w, "shell", gameSession); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while executing template: %v\n", err)
		os.Exit(1)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		pattern = "penta_decathlon"
	}
	game, err := spawnGame(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to spawn a new game: %v\n", err)
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
				gameSession.Game.LoadState(filepath.Join(savePath, gameSession.Id+".json"), 0)
			case "saveState":
				gameSession.Game.SaveState(filepath.Join(savePath, gameSession.Id+".json"), 0)
			case "togglePause":
				gameSession.Game.Running = !gameSession.Game.Running
				ioSocket.MessageChannel <- io.SocketMessage{Action: "pauseToggled", Payload: gameSession.Game.Running}
			case "step":
				if gameSession.Game.Running {
					gameSession.Game.Running = false
					ioSocket.MessageChannel <- io.SocketMessage{Action: "pauseToggled", Payload: gameSession.Game.Running}
				}
				gameSession.Game.Step()
			case "close":
				gameSession.SignalClose(&ioSocket)
			}
		}
	}

	gameSession.ConnectIO(&ioSocket, listenerChannel, eventHandler)

	go ioSocket.ListenEvents()
	go ioSocket.SendMessages()
	// TODO: Delete game save file on last client's connection end
}

func main() {
	defer func() {
		// TODO: Check for leaked resources, channels, websockets, gamesession and multiplexer
		// TODO: Implement if required
	}()

	staticFs := http.FileServer(http.Dir(filepath.Join(os.Getenv("GOPATH"), "web", "static")))

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
