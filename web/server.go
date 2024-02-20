package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/patterns"
)

var cellMap cell_map.Map

func getServerDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return filepath.Dir(filename)
}

func gameViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join(getServerDir(), "templates", "index.tmpl"),
		filepath.Join(getServerDir(), "templates", "cellMap.tmpl"),
	))
	tmpl.ExecuteTemplate(w, "index", cellMap.Rasterize(10))
}

func nextStepHandler(w http.ResponseWriter, r *http.Request) {
	cellMap.Step()
	tmpl := template.Must(template.ParseFiles(
		filepath.Join(getServerDir(), "templates", "cellMap.tmpl"),
	))
	tmpl.ExecuteTemplate(w, "cellMap", cellMap.Rasterize(10))
}

func main() {
	cellMap = cell_map.Create()
	cells := patterns.GetPrimitive("PentaDecathlon", 0, 0)
	cellMap.AddCells(cells)

	staticFs := http.FileServer(http.Dir(filepath.Join(getServerDir(), "static")))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.Handle("/game/", http.HandlerFunc(gameViewHandler))
	mux.Handle("/step/", http.HandlerFunc(nextStepHandler))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}