package tui

import (
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

	for {
		if app.OnEvent(termbox.PollEvent()) {
			break
		}

		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		drawWidgets(widgets)
		termbox.Flush()
	}

	return nil
}
