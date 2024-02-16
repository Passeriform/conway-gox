package io

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

// TODO: Implement inheritance from IO base class
type Terminal struct {
	screen    tcell.Screen
	width     int
	height    int
	zoomLevel float64
}

var aliveCell = "\u2588"
var deadCell = "\u2591"

func emitStr(s tcell.Screen, x int, y int, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, tcell.StyleDefault)
		x += w
	}
}

func Create(w int, h int, z float64) Terminal {
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

	return Terminal{s, w, h, z}
}

func (t *Terminal) Blit(cell_map cell_map.Map) {
	bounds := cell_map.GetBounds()

	blitMap := make([][]bool, max(bounds.Bottom-bounds.Top+1, (int)((float64)(t.height)/t.zoomLevel)))

	for r := range blitMap {
		blitMap[r] = make([]bool, max(bounds.Right-bounds.Left+1, (int)((float64)(t.height)/t.zoomLevel)))
	}

	for _, cell := range cell_map.GetCells() {
		row, column := cell.GetPosition()
		blitMap[row+(len(blitMap)/2)][column+(len(blitMap[0])/2)] = true
	}

	t.screen.Clear()
	for r := range blitMap {
		for c, hasCell := range blitMap[r] {
			if hasCell {
				emitStr(t.screen, r, c, aliveCell)
			} else {
				emitStr(t.screen, r, c, deadCell)
			}
		}
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
