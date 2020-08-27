package tui

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

// Surround coordinates with box
func DrawBox(x, y, w, h int) {
	const coldef = termbox.ColorDefault

	if runewidth.EastAsianWidth {
		termbox.SetCell(x-1, y-1, '+', coldef, coldef)
		termbox.SetCell(x+w, y-1, '+', coldef, coldef)
		termbox.SetCell(x+w, y+h, '+', coldef, coldef)
		termbox.SetCell(x-1, y+h, '+', coldef, coldef)
		fill(x, y-1, w, 1, termbox.Cell{Ch: '-'})
		fill(x, y+h, w, 1, termbox.Cell{Ch: '-'})
		fill(x-1, y, 1, h, termbox.Cell{Ch: '|'})
		fill(x+w, y, 1, h, termbox.Cell{Ch: '|'})
	} else {
		termbox.SetCell(x-1, y-1, '┌', coldef, coldef)
		termbox.SetCell(x+w, y-1, '┐', coldef, coldef)
		termbox.SetCell(x+w, y+h, '┘', coldef, coldef)
		termbox.SetCell(x-1, y+h, '└', coldef, coldef)
		fill(x, y-1, w, 1, termbox.Cell{Ch: '─'})
		fill(x, y+h, w, 1, termbox.Cell{Ch: '─'})
		fill(x-1, y, 1, h, termbox.Cell{Ch: '│'})
		fill(x+w, y, 1, h, termbox.Cell{Ch: '│'})
	}
}

type Widget interface {
	Draw()
}
