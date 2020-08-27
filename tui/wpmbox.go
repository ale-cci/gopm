package tui

import (
	"github.com/nsf/termbox-go"
	"strings"
)

type WpmBox struct {
	buffer [][]termbox.Cell
	text   []rune
	x, y   int
	w, h   int
	// current line of text and cursor position relative to the line
	line, cursor   int
	wrong, correct int
}

const (
	CELL_WRONG = 0
	CELL_RIGHT = 1
	CELL_BLANK = 2
)

func NewWpmBox(x int, y int, w int, h int, text []rune) *WpmBox {
	box := &WpmBox{x: x, y: y, w: w, h: h, text: text}
	box.SetText(text)
	return box
}

func (w *WpmBox) SetText(text []rune) {
	w.buffer = [][]termbox.Cell{}
	lines := strings.Split(string(w.text), "\n")
	lastLine := len(lines) - 1

	for curLine, str := range lines {
		cells := make([]termbox.Cell, len(str))

		for i, char := range str {
			if char == ' ' {
				char = '·'
			}
			cells[i] = termbox.Cell{Ch: char, Fg: termbox.ColorDefault, Bg: termbox.ColorDefault}
		}

		if lastLine != curLine {
			cells = append(cells, termbox.Cell{Ch: '↩', Fg: termbox.ColorDefault, Bg: termbox.ColorDefault})
		}
		w.buffer = append(w.buffer, cells)
	}
}

func (w *WpmBox) cellColor(currentChar int) (termbox.Attribute, termbox.Attribute) {
	switch w.CellType(currentChar) {
	case CELL_RIGHT:
		return termbox.ColorGreen, termbox.ColorDefault
	case CELL_WRONG:
		return termbox.ColorMagenta, termbox.ColorDefault
	default:
		return termbox.ColorDefault, termbox.ColorDefault
	}
}

func (w *WpmBox) CellType(currentChar int) int {
	switch {
	case currentChar < w.correct:
		return CELL_RIGHT
	case currentChar >= w.correct && currentChar < w.correct+w.wrong:
		return CELL_WRONG
	default:
		return CELL_BLANK
	}
}

func (w *WpmBox) Draw() {
	DrawBox(w.x, w.y, w.w, w.h)

	// Draw box content
	currentChar := 0
	for y, row := range w.buffer {
		for x, cell := range row {
			// Determine color based on correctness
			colFg, colBg := w.cellColor(currentChar)
			termbox.SetCell(w.x+x, w.y+y, cell.Ch, colFg, colBg)
			currentChar++
		}
	}

	// Draw cursor
	termbox.SetCursor(w.x+w.cursor, w.y+w.line)
}

func (w *WpmBox) InsKey(key rune) {
	if w.TextEnded() {
		return
	}

	w.incCursor()

	if w.CurrentRune() != key || w.wrong >= 1 {
		w.wrong += 1
	} else {
		w.correct += 1
	}
}

func (w *WpmBox) Backspace() {
	if w.wrong > 0 {
		w.wrong -= 1
	} else if w.correct > 0 {
		w.correct -= 1
	}
	w.decCursor()
}

func (w *WpmBox) TextEnded() bool {
	idx := w.correct + w.wrong
	return idx == len(w.text)
}

func (w *WpmBox) CurrentRune() rune {
	idx := w.correct + w.wrong
	return w.text[idx]
}

func (w *WpmBox) incCursor() {
	if w.CurrentRune() == '\n' {
		w.line++
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
		w.cursor = len(w.buffer[w.line]) - 1
	}
}
