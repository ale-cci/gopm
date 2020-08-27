package tui

import (
	"github.com/ale-cci/gopm/quotes"
	"github.com/nsf/termbox-go"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type WpmBox struct {
	x, y int
	w, h int
	// current line of text and cursor position relative to the line
	textStructure []int

	line, cursor int
	CurrentText  quotes.CurrentText

	ScrollOff int
	offset    int
}

// Possible text character statuses
const (
	CELL_WRONG = 0
	CELL_RIGHT = 1
	CELL_BLANK = 2
)

// WpmBox constructor
func NewWpmBox(x int, y int, w int, h int, text string) *WpmBox {
	box := &WpmBox{x: x, y: y, w: w, h: h}
	box.SetText(text)
	box.offset = 0
	return box
}

// Set wpmbox text, update internal buffer
func (w *WpmBox) SetText(text string) {
	// Reset cursor position
	w.cursor = 0
	w.line = 0

	w.CurrentText = quotes.CurrentText{Text: text}

	lines := strings.Split(text, "\n")
	w.textStructure = make([]int, len(lines))

	for i, line := range lines {
		w.textStructure[i] = len(line)

		// new line character to all lines except last one
		if i != len(lines)-1 {
			w.textStructure[i] += 1
		}
	}
}

func (w *WpmBox) RuneAt(position int) rune {
	char := w.CurrentText.RuneAt(position)
	switch char {
	case ' ':
		return '·'
	case '\n':
		return '↩'
	case '\t':
		return '⇥'
	default:
		return char
	}
}

// Color of text character at given position
func (w *WpmBox) cellColor(currentChar int) (termbox.Attribute, termbox.Attribute) {

	switch w.CurrentText.CharStatus(currentChar) {
	case quotes.RIGHT:
		return termbox.ColorGreen, termbox.ColorDefault
	case quotes.WRONG:
		return termbox.ColorDefault, termbox.ColorRed
	default:
		return termbox.ColorDefault, termbox.ColorDefault
	}
}

func (w *WpmBox) Draw() {
	DrawBox(w.x, w.y, w.w, w.h)

	// Draw box content
	currentChar := 0
	for l := 0; l < w.offset; l++ {
		currentChar += w.textStructure[l]
	}

	for y := 0; y < min(w.h, len(w.textStructure)-w.offset); y++ {
		for x := 0; x < w.textStructure[w.offset+y]; x++ {
			// Determine color based on correctness
			fg, bg := w.cellColor(currentChar)
			char := w.RuneAt(currentChar)

			termbox.SetCell(w.x+x, w.y+y, char, fg, bg)
			currentChar++
		}
	}

	// Draw cursor
	termbox.SetCursor(w.x+w.cursor, w.y+w.line)
}

func (w *WpmBox) InsKey(key rune) {
	if w.CurrentText.IsEndPosition() {
		return
	}

	w.incCursor()
	w.CurrentText.InsKey(key)
}

func (w *WpmBox) Backspace() {
	w.CurrentText.Backspace()
	w.decCursor()
}

func (w *WpmBox) incCursor() {
	if w.CurrentText.CurrentRune() == '\n' {
		if w.h-w.line > w.ScrollOff {
			w.line++
		} else {
			w.offset++
		}

		w.cursor = 0
	} else {
		w.cursor++
	}
}

func (w *WpmBox) decCursor() {
	if w.cursor > 0 {
		w.cursor--
	} else if w.line > 0 {
		w.line--
		w.cursor = w.textStructure[w.line] - 1
	}
}
