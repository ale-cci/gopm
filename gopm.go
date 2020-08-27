package main

import (
	"flag"
	"fmt"
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
	app.box.ScrollOff = 4

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
	filename := flag.String("file", "", "test")
	flag.Parse()

	if *filename == "" {
		panic("Flag --file not provided")
	}
	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		fmt.Printf("Unable to open file %q\n", *filename)
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	quote, _ := quotes.LoadFile(file)
	quoteList := []quotes.Quote{*quote}

	qi := quotes.NewQuoteIterator(quoteList)
	app := NewApp(qi)

	tui.Run(app)
}
