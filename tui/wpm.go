package tui

import (
	"fmt"
	"time"

	"github.com/ale-cci/gopm/wpm"
	"github.com/nsf/termbox-go"
)

type WpmBar struct {
	x, y    int
	w, h    int
	Counter *wpm.KeystrokeCounter
}

func NewWpmBar(x, y, w, h int, counter *wpm.KeystrokeCounter) *WpmBar {
	return &WpmBar{Counter: counter, x: x, y: y, w: w, h: h}
}
func (w *WpmBar) Draw() {
	DrawBox(w.x, w.y, w.w, w.h)
	score := w.Counter.Wpm(time.Now())

	runes := []rune(fmt.Sprintf("Wpm: %5.1f | Accuracy: TODO", score))
	const colordef = termbox.ColorDefault

	for x, char := range runes {
		termbox.SetCell(w.x+x, w.y, char, colordef, colordef)
	}
}
