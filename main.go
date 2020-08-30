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

	app.box = tui.NewWpmBox(1, 3, maxX-2, maxY-4, app.text, &app.counter)
	app.reloadText()
	app.box.ScrollOff = 4

	return []tui.Widget{
		tui.NewWpmBar(1, 1, maxX-2, 1, &app.counter),
		app.box,
	}
}

func (app *App) reloadText() {
	text := app.iterator.Current()
	app.counter.Reset()
	app.counter.Capacity = len(text)

	app.text = tui.NewParsedText(text)
	app.box.SetText(app.text)
}

func (app *App) insKey(char rune) {
	if app.counter.IsStartPosition() {
		app.counter.Start = time.Now()
	}

	correct := app.text.RuneAt(app.counter.Inserted()) == char
	app.box.IncCursor()
	app.counter.InsKey(correct)
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
			app.reloadText()

		case termbox.KeyCtrlP:
			app.iterator.Prev()
			app.reloadText()

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			app.box.DecCursor()
			app.counter.Backspace()

		default:
			if ev.Ch != 0 {
				app.insKey(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	}

	if app.counter.IsEndPosition() {
		app.iterator.Next()
		app.reloadText()
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
