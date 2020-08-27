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
	switch {
	case currentChar < w.correct:
		return termbox.ColorGreen, termbox.ColorDefault
	case currentChar >= w.correct && currentChar < w.wrong:
		return termbox.ColorRed, termbox.ColorDefault
	default:
		return termbox.ColorDefault, termbox.ColorDefault
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
	idx := w.correct + w.wrong

	if w.text[idx] != key || w.wrong >= 1 {
		w.wrong += 1
	} else {
		w.correct += 1
	}
	w.incCursor()
}

func (w *WpmBox) Backspace() {
	if w.wrong > 0 {
		w.wrong -= 1
	} else if w.correct > 0 {
		w.correct -= 1
	}
	w.decCursor()
}

func (w *WpmBox) incCursor() {
}

func (w *WpmBox) decCursor() {
}
