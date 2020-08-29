package tui

import (
	"github.com/ale-cci/gopm/quotes"
	"github.com/nsf/termbox-go"
)

func RuneMap() map[rune]rune {
	return map[rune]rune{
		' ':  '·',
		'\n': '↩',
		'\t': '⇥',
	}
}

func ColorMap() map[quotes.CharStatus]CellColor {
	return map[quotes.CharStatus]CellColor{
		quotes.RIGHT: {Fg: termbox.ColorGreen, Bg: termbox.ColorDefault},
		quotes.WRONG: {Fg: termbox.ColorDefault, Bg: termbox.ColorRed},
		quotes.BLANK: {Fg: termbox.ColorDefault, Bg: termbox.ColorDefault},
	}
}

func SpecialRuneSizes() map[rune]int {
	return map[rune]int{
		'\t': 4,
	}
}
