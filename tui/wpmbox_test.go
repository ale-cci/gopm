package tui

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func asStr(line []termbox.Cell) string {
	var str string
	for _, cell := range line {
		str += string(cell.Ch)
	}
	return str
}

func TestWpmBox(t *testing.T) {

	t.Run("Should correctly initialize buffer", func(t *testing.T) {
		text := []rune(`Test multiline string
second line here`)
		box := NewWpmBox(0, 0, 0, 0, text)

		if len(box.buffer) != 2 {
			t.Errorf("Different number of line parsed: %d, expected 2", len(box.buffer))
			return
		}

		firstLine := asStr(box.buffer[0])
		expected := "Test·multiline·string↩"

		if firstLine != expected {
			t.Errorf("String parsed incorrectly: %s, expected %s", firstLine, expected)
			return
		}

		secondLine := asStr(box.buffer[1])
		expected = "second·line·here"

		if secondLine != expected {
			t.Errorf("String parsed incorrectly: %s, expected %s", secondLine, expected)
			return
		}
	})

	t.Run("Should clear previous text", func(t *testing.T) {
		text := []rune("First string")
		box := NewWpmBox(0, 0, 0, 0, text)
		box.SetText([]rune("Second string"))

		if len(box.buffer) != 1 {
			t.Errorf("Multiple lines (%d) found in box, expected 1", len(box.buffer))
			return
		}
	})

}
