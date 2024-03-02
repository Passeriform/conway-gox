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
	screen          tcell.Screen
	MessageChannel  chan cell_map.Map
	listenerChannel chan tcell.Event
	once            sync.Once
}

var aliveCell rune = '\u2B1C'

func NewTerminal() (Terminal, <-chan tcell.Event, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, err := tcell.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create an instance of the terminal screen: %v\n", err)
		return Terminal{}, nil, err
	}

	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize the terminal screen: %v\n", err)
		return Terminal{}, nil, err
	}

	s.SetStyle(tcell.StyleDefault)

	s.Clear()

	lChan := make(chan tcell.Event)
	mChan := make(chan cell_map.Map)

	return Terminal{screen: s, listenerChannel: lChan, MessageChannel: mChan}, lChan, nil
}

func (t *Terminal) SendMessages() {
	defer func() {
		t.Close()
	}()

	for message := range t.MessageChannel {
		t.screen.Clear()
		width, height := t.screen.Size()
		for _, cell := range message.EncodeJson(0) {
			t.screen.SetContent(
				(width/2)+(2*cell[1]),
				(height/2)+(2*cell[0]),
				aliveCell,
				nil,
				tcell.StyleDefault,
			)
		}
		t.screen.Show()
	}
}

func (t *Terminal) Blit(mapChannel <-chan cell_map.Map) {
	defer func() {
		t.Close()
	}()

	for cm := range mapChannel {
		t.MessageChannel <- cm
	}
}

func (t *Terminal) ListenEvents() {
	defer func() {
		t.Close()
	}()

	for {
		ev := t.screen.PollEvent()
		t.listenerChannel <- ev
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
