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

	t.Run("Cursor should be positioned on origin at the start", func(t *testing.T) {
		text := []rune("Example text")
		box := NewWpmBox(0, 0, 0, 0, text)

		if box.line != 0 || box.cursor != 0 {
			t.Errorf("Cursor not in origin position: (%d, %d) expected (0, 0)", box.cursor, box.line)
		}
	})

	t.Run("Cursor should normally advance horizontally", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, []rune("Example"))
		box.incCursor()

		if box.line != 0 || box.cursor != 1 {
			t.Errorf("Wrong cursor position: (%d, %d) expected (1, 0)", box.cursor, box.line)
		}
	})

	t.Run("incCursor", func(t *testing.T) {
		t.Run("Cursor should follow new line", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("\n"))
			box.incCursor()
			if box.line != 1 || box.cursor != 0 {
				t.Errorf("Wrong cursor position: (%d, %d) expected (0, 1)", box.cursor, box.line)
			}
		})

		t.Run("Cursor should follow new line when line is ended", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test\nline"))
			box.correct = 4
			box.cursor = 4

			box.incCursor()

			expectedX := 0
			expectedY := 1
			if box.line != expectedY || box.cursor != expectedX {
				t.Errorf("Wrong cursor position: (%d, %d) expected (%d, %d)", box.cursor, box.line, expectedX, expectedY)
			}
		})
		t.Run("Should not go to new line before last char", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test \nline"))
			box.correct = 4
			box.cursor = 4

			box.incCursor()

			expectedX := 5
			expectedY := 0
			if box.line != expectedY || box.cursor != expectedX {
				t.Errorf("Wrong cursor position: (%d, %d) expected (%d, %d)", box.cursor, box.line, expectedX, expectedY)
			}
		})

		t.Run("Should not increment function when text is ended", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test"))
			box.correct = 4
			box.cursor = 4

			box.InsKey(' ')
			if box.wrong != 0 || box.correct != 4 {
				t.Errorf("Incremented letters counted; wrong: %d, correct: %d, expected: (0, 4)", box.wrong, box.correct)
				return
			}

			if box.cursor != 4 || box.line != 0 {
				t.Errorf("Changed cursor position: (%d, %d), expected (4, 0)", box.cursor, box.line)
				return
			}
		})

		t.Run("InsKey behaviour on newline inserted", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test\na"))
			box.correct = 3 // tesT
			box.cursor = 3

			box.InsKey('t')
			if box.correct != 4 {
				t.Errorf("Wrong number of correct letters: %d, expected 4", box.correct)
			}

			if box.line != 0 || box.cursor != 4 {
				t.Errorf("Wrong cursor position: (%d, %d) expected (0, 4)", box.line, box.cursor)
			}

		})
	})

	t.Run("decCursor", func(t *testing.T) {
		t.Run("Should normally decremnt cursor position on decCursor called", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test"))
			box.cursor = 2
			box.correct = 2

			box.decCursor()
			if box.cursor != 1 {
				t.Errorf("Unexpected cursor position: %d, expected 1", box.cursor)
			}
		})

		t.Run("Should not decrement cursor at beginning of text", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test"))

			box.decCursor()
			if box.cursor != 0 {
				t.Errorf("Unexpected cursor position: %d, expected 1", box.cursor)
			}
		})

		t.Run("Should return to previous line on backpress", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, []rune("test\n"))
			box.correct = 4
			box.cursor = 0
			box.line = 1

			box.decCursor()
			if box.cursor != 4 || box.line != 0 {
				t.Errorf("Unexpected cursor position: (%d, %d), expected (4, 0)", box.cursor, box.line)
			}
		})
	})

	t.Run("Should identify wrong cells", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, []rune("text"))
		box.correct = 2
		box.wrong = 1

		if box.CellType(0) != CELL_RIGHT {
			t.Errorf("Failed to identify right cell")
		}

		if box.CellType(2) != CELL_WRONG {
			t.Errorf("Failed to identify wrong cell")
		}

		if box.CellType(3) != CELL_BLANK {
			t.Errorf("Failed to identify blank cell")
		}
	})
}
