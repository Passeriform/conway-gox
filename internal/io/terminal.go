package io

import (
	"fmt"
	"os"
	"sync"

	"github.com/gdamore/tcell"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

// TODO: Implement generic interface for all IO handlers and use in GameSession
type Terminal struct {
	screen    tcell.Screen
	Listeners []chan<- tcell.Event
	zoomLevel float64
	once      sync.Once
}

var aliveCell rune = '\u2B1C'

func NewTerminal() (Terminal, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, err := tcell.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create an instance of the terminal screen: %v\n", err)
		return Terminal{}, err
	}

	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize the terminal screen: %v\n", err)
		return Terminal{}, err
	}

	s.SetStyle(tcell.StyleDefault)

	s.Clear()

	return Terminal{screen: s, zoomLevel: 1}, nil
}

func (t *Terminal) Blit(mapChannel <-chan cell_map.Map) {
	defer func() {
		t.Close()
	}()

	for m := range mapChannel {
		t.screen.Clear()
		width, _ := t.screen.Size()
		for _, cell := range m.EncodeJson(width / 2) {
			t.screen.SetContent(cell[1], cell[0], aliveCell, nil, tcell.StyleDefault)
		}
		t.screen.Show()
	}
}

func (s *Terminal) AddListener(eventChannel chan<- tcell.Event) {
	s.Listeners = append(s.Listeners, eventChannel)
}

func (t *Terminal) ListenEvents() {
	defer func() {
		t.Close()
	}()

	for {
		ev := t.screen.PollEvent()
		for _, listener := range t.Listeners {
			listener <- ev
		}
		switch ev := ev.(type) {
		case *tcell.EventResize:
			t.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}
}

func (t *Terminal) Close() {
	t.once.Do(func() {
		t.screen.Fini()
	})
}
