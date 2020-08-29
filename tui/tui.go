package tui

import (
	"time"

	"github.com/nsf/termbox-go"
)

type TUIApp interface {
	Build(maxX int, maxY int) []Widget
	OnEvent(e termbox.Event) bool
}

// Render all widgets on screen
func drawWidgets(widgets []Widget) {
	for _, widget := range widgets {
		widget.Draw()
	}
}

func Run(app TUIApp) error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	maxX, maxY := termbox.Size()
	widgets := app.Build(maxX, maxY)

	drawWidgets(widgets)
	termbox.Flush()

	eventQueue := make(chan termbox.Event)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	drawTick := time.NewTicker(30 * time.Millisecond)

mainLoop:
	for {

		select {
		case ev := <-eventQueue:
			if app.OnEvent(ev) {
				break mainLoop
			}
		case <-drawTick.C:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawWidgets(widgets)
			termbox.Flush()
		}
	}

	return nil
}
