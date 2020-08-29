package tui

import (
	"github.com/ale-cci/gopm/wpm"
	"github.com/nsf/termbox-go"
)

func RuneMap() map[rune]rune {
	return map[rune]rune{
		' ':  '·',
		'\n': '↩',
		'\t': '⇥',
	}
}

func ColorMap() map[wpm.CharStatus]CellColor {
	return map[wpm.CharStatus]CellColor{
		wpm.RIGHT: {Fg: termbox.ColorGreen, Bg: termbox.ColorDefault},
		wpm.WRONG: {Fg: termbox.ColorDefault, Bg: termbox.ColorRed},
		wpm.BLANK: {Fg: termbox.ColorDefault, Bg: termbox.ColorDefault},
	}
}

func SpecialRuneSizes() map[rune]int {
	return map[rune]int{
		'\t': 4,
	}
}
