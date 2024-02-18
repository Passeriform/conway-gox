package io

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

// TODO: Implement inheritance from IO base class
type Terminal struct {
	screen    tcell.Screen
	zoomLevel float64
}

var aliveCell rune = '\u2B1C'
var deadCell rune = '\u2B1B'

func Create() Terminal {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, err := tcell.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault)

	s.Clear()

	return Terminal{s, 1}
}

func (t *Terminal) Blit(cell_map cell_map.Map) {
	t.screen.Clear()
	width, height := t.screen.Size()
	widthOffset, heightOffset := width/2, height/2
	for _, cell := range cell_map.GetCells() {
		row, column := cell.GetPosition()
		t.screen.SetContent(column+widthOffset, row+heightOffset, aliveCell, nil, tcell.StyleDefault)
	}
	t.screen.Show()
}

func (t *Terminal) ListenEvents(event chan<- tcell.Event) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	for {
		ev := t.screen.PollEvent()
		event <- ev
		switch ev := ev.(type) {
		case *tcell.EventResize:
			t.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				t.Close()
				return true
			}
		}
	}
}

func (t *Terminal) Close() {
	t.screen.Fini()
}
