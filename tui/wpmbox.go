package tui

import (
	"strings"

	"github.com/ale-cci/gopm/quotes"
	"github.com/nsf/termbox-go"
)

type WpmBox struct {
	// Box coordintates
	x, y int
	w, h int

	// Effective length of each line of the text
	textStructure []int

	// current line and position of cursor on the screen
	line, cursor     int
	KeystrokeCounter quotes.KeystrokeCounter

	// When number of line from margin where screen scrolling should start
	ScrollOff int

	// Number of lines hidden above
	offset int
}

// WpmBox constructor
func NewWpmBox(x int, y int, w int, h int, text string) *WpmBox {
	box := &WpmBox{x: x, y: y, w: w, h: h}
	box.SetText(text)
	return box
}

// Set wpmbox text, update internal buffer
func (w *WpmBox) SetText(text string) {
	// Reset cursor position
	w.cursor = 0
	w.line = 0
	w.offset = 0

	w.KeystrokeCounter = quotes.KeystrokeCounter{Text: text}

	lines := strings.Split(text, "\n")
	w.textStructure = make([]int, len(lines))

	for i, line := range lines {
		// w.textStructure[i] = len(line)
		for _, char := range line {
			w.textStructure[i] += runeSize(char)
		}

		// new line character to all lines except last one
		if i != len(lines)-1 {
			w.textStructure[i] += 1
		}
	}
}

func runeSize(char rune) int {
	val, ok := SpecialRuneSizes()[char]
	if ok {
		return val
	}
	return 1
}

func (w *WpmBox) cellColor(currentChar int) CellColor {
	val, ok := ColorMap()[w.KeystrokeCounter.CharStatus(currentChar)]

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
		currentChar += w.textStructure[l]
	}

	// Draw box content
	for y := 0; y < min(w.h, len(w.textStructure)-w.offset); y++ {
		for x := 0; x < w.textStructure[w.offset+y]; {
			color := w.cellColor(currentChar)

			currChar := w.KeystrokeCounter.RuneAt(currentChar)
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

func (w *WpmBox) InsKey(key rune) {
	if w.KeystrokeCounter.IsEndPosition() {
		return
	}

	w.incCursor()
	w.KeystrokeCounter.InsKey(key)
}

func (w *WpmBox) Backspace() {
	w.KeystrokeCounter.Backspace()
	w.decCursor()
}

func (w *WpmBox) incCursor() {
	if w.KeystrokeCounter.CurrentRune() == '\n' {
		if w.h-w.line > w.ScrollOff {
			w.line++
		} else if w.line+w.offset+w.ScrollOff < len(w.textStructure)-1 {
			w.offset++
		} else {
			w.line++
		}

		w.cursor = 0
	} else {
		w.cursor += runeSize(w.KeystrokeCounter.CurrentRune())
	}
}

func (w *WpmBox) decCursor() {
	c := w.KeystrokeCounter.CurrentRune()
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
		w.cursor = w.textStructure[w.line+w.offset] - 1
	}
}
