package tui

import (
	"github.com/ale-cci/gopm/wpm"
	"github.com/nsf/termbox-go"
)

type WpmBox struct {
	// Box coordintates
	x, y int
	w, h int

	// Effective length of each line of the Text
	counter *wpm.KeystrokeCounter
	pt      *ParsedText

	// current line and position of cursor on the screen
	line, cursor int

	// When number of line from margin where screen scrolling should start
	ScrollOff int

	// Number of lines hidden above
	offset int
}

// WpmBox constructor
func NewWpmBox(x int, y int, w int, h int, pt *ParsedText, counter *wpm.KeystrokeCounter) *WpmBox {
	box := &WpmBox{x: x, y: y, w: w, h: h, pt: pt, counter: counter}
	box.Reset()
	return box
}

func (wb *WpmBox) SetText(pt *ParsedText) {
}

// Set wpmbox Text, update internal buffer
func (w *WpmBox) Reset() {
	// Reset cursor position
	w.cursor = 0
	w.line = 0
	w.offset = 0
}

func runeSize(char rune) int {
	val, ok := SpecialRuneSizes()[char]
	if ok {
		return val
	}
	return 1
}

func (w *WpmBox) cellColor(currentChar int) CellColor {
	val, ok := ColorMap()[w.counter.CharStatus(currentChar)]

	if !ok {
		panic("Invalid cell state")
	}
	return val
}

func parseRune(char rune) rune {
	val, ok := RuneMap()[char]

	if ok {
		return val
	}
	return char
}

func (w *WpmBox) Draw() {
	// DrawBox(w.x, w.y, w.w, w.h)

	currentChar := 0
	for l := 0; l < w.offset; l++ {
		currentChar += w.pt.Structure[l]
	}

	// Draw box content
	for y := 0; y < min(w.h, len(w.pt.Structure)-w.offset); y++ {
		for x := 0; x < w.pt.Structure[w.offset+y]; {
			color := w.cellColor(currentChar)

			currChar := w.pt.RuneAt(currentChar)
			parsedChar := parseRune(currChar)

			termbox.SetCell(w.x+x, w.y+y, parsedChar, color.Fg, color.Bg)
			currentChar++

			size := runeSize(currChar)
			x += size
		}
	}

	// Draw cursor
	termbox.SetCursor(w.x+w.cursor, w.y+w.line)
}
func (wb *WpmBox) currentRune() rune {
	return wb.pt.RuneAt(wb.counter.Position())
}

func (w *WpmBox) IncCursor() {
	if w.counter.IsEndPosition() {
		return
	}
	prevRune := w.pt.RuneAt(w.counter.Position() - 1)
	if prevRune == '\n' {
		if w.h-w.line > w.ScrollOff {
			w.line++
		} else if w.line+w.offset+w.ScrollOff < len(w.pt.Structure)-1 {
			w.offset++
		} else {
			w.line++
		}

		w.cursor = 0
	} else {
		w.cursor += runeSize(prevRune)
	}
}

func (w *WpmBox) DecCursor() {
	c := w.currentRune()
	if w.cursor > 0 {
		w.cursor -= runeSize(c)
	} else if w.line > 0 {
		if w.line-w.offset > w.ScrollOff {
			w.line--
		} else if w.offset > 0 {
			w.offset--
		} else {
			w.line--
		}
		w.cursor = w.pt.Structure[w.line+w.offset] - 1
	}
}
