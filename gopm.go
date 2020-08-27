package main

import (
	"github.com/nsf/termbox-go"
	"os"

	"github.com/ale-cci/gopm/quotes"
	"github.com/ale-cci/gopm/tui"
)

type App struct {
	qi  *quotes.QuoteIterator
	box *tui.WpmBox
}

func (app *App) Build(maxX int, maxY int) []tui.Widget {
	text := app.qi.Current().Text
	app.box = tui.NewWpmBox(1, 1, maxX-2, maxY-2, text)

	return []tui.Widget{app.box}
}

func (app *App) OnEvent(ev termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyCtrlC:
			return true
		case termbox.KeySpace:
			app.box.InsKey(' ')
		case termbox.KeyEnter:
			app.box.InsKey('\n')
		case termbox.KeyTab:
			app.box.InsKey('\t')

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			app.box.Backspace()
		default:
			if ev.Ch != 0 {
				app.box.InsKey(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	}

	if app.box.CurrentText.IsEndPosition() {
		app.qi.Next()
		app.box.SetText(app.qi.Current().Text)
	}
	return false
}

func NewApp(qi *quotes.QuoteIterator) *App {
	return &App{qi: qi}
}

func main() {
	file, err := os.Open("gopm.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	quote, _ := quotes.LoadFile(file)
	quoteList := []quotes.Quote{*quote}

	/*
		file, err := os.Open("default.json")
		defer file.Close()

		if err != nil {
			panic(err)
		}
		quoteList, err := quotes.LoadFromJson(file)
		if err != nil {
			panic(err)
		}
	*/

	qi := quotes.NewQuoteIterator(quoteList)
	app := NewApp(qi)

	tui.Run(app)
}
