package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/ale-cci/gopm/chunk"
	"github.com/ale-cci/gopm/tui"
	"github.com/ale-cci/gopm/wpm"
)

type App struct {
	box      *tui.WpmBox
	iterator chunk.Iterator
	counter  wpm.KeystrokeCounter
	text     *tui.ParsedText
}

func (app *App) Build(maxX int, maxY int) []tui.Widget {
	app.counter = wpm.KeystrokeCounter{}

	text := app.iterator.Current()
	app.setText(text)

	app.box = tui.NewWpmBox(1, 3, maxX-2, maxY-4, app.text, &app.counter)
	app.box.ScrollOff = 4

	return []tui.Widget{
		tui.NewWpmBar(1, 1, maxX-2, 1, &app.counter),
		app.box,
	}
}

func (app *App) setText(text string) {
	app.text = tui.NewParsedText(text)
	app.box.SetText(&app.text)
	app.counter.Capacity = len(text)
	app.counter.Reset()
}

func (app *App) insKey(char rune) {
	correct := app.text.RuneAt(app.counter.Position()) == char
	app.counter.InsKey(correct)
	app.box.IncCursor()
}

func (app *App) OnEvent(ev termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyCtrlC:
			return true
		case termbox.KeySpace:
			app.insKey(' ')

		case termbox.KeyEnter:
			app.insKey('\n')
		case termbox.KeyTab:
			app.insKey('\t')

		case termbox.KeyCtrlN:
			app.iterator.Next()
			app.setText(app.iterator.Current())

		case termbox.KeyCtrlP:
			app.iterator.Prev()
			app.setText(app.iterator.Current())

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			app.counter.Backspace()
			app.box.DecCursor()

		default:
			if ev.Ch != 0 {
				if app.counter.IsStartPosition() {
					app.counter.Start = time.Now()
				}
				app.insKey(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	}

	if app.counter.IsEndPosition() {
		app.iterator.Next()
		app.setText(app.iterator.Current())
	}
	return false
}

func NewApp(chunkIterator chunk.Iterator) *App {
	return &App{iterator: chunkIterator}
}

func main() {
	flag.Parse()
	filenames := flag.Args()
	if len(filenames) == 0 {
		fmt.Println("No filename provided")
		return
	}

	chunkedFiles := make([]chunk.Iterator, len(filenames))

	for i, filename := range filenames {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("Unable to open file %q\n", filename)
			os.Exit(1)
		}
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		chunkedFiles[i] = chunk.NewChunkedFile(file, 40)
	}

	chunkIterator := &chunk.ChunkIterator{Files: chunkedFiles}
	app := NewApp(chunkIterator)

	tui.Run(app)
	return
}
