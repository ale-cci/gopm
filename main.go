package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/ale-cci/gopm/chunk"
	"github.com/ale-cci/gopm/tui"
)

type App struct {
	box           *tui.WpmBox
	chunkIterator *chunk.ChunkIterator
}

func (app *App) Build(maxX int, maxY int) []tui.Widget {
	text := app.chunkIterator.Current()
	app.box = tui.NewWpmBox(1, 3, maxX-2, maxY-4, text)
	app.box.ScrollOff = 4

	return []tui.Widget{
		tui.NewWpmBar(1, 1, maxX-2, 1, &app.box.KeystrokeCounter),
		app.box,
	}
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

		case termbox.KeyCtrlN:
			app.chunkIterator.Next()
			app.box.SetText(app.chunkIterator.Current())

		case termbox.KeyCtrlP:
			app.chunkIterator.Prev()
			app.box.SetText(app.chunkIterator.Current())

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			app.box.Backspace()

		default:
			if ev.Ch != 0 {
				if app.box.KeystrokeCounter.IsStartPosition() {
					app.box.KeystrokeCounter.Start = time.Now()
				}
				app.box.InsKey(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	}

	if app.box.KeystrokeCounter.IsEndPosition() {
		app.chunkIterator.Next()
		app.box.SetText(app.chunkIterator.Current())
	}
	return false
}

func NewApp(chunkIterator *chunk.ChunkIterator) *App {
	return &App{chunkIterator: chunkIterator}
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
