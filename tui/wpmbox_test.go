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
		text := "multiline\nstring"
		box := NewWpmBox(0, 0, 0, 0, text)

		got := len(box.textStructure)
		expected := 2

		if got != expected {
			t.Errorf("Different line of text parsed: %d, expected: %d", got, expected)
		} else {

			got = box.textStructure[0]
			expected = len("multiline\n")

			if got != expected {
				t.Errorf("Mismatched number of characters in first line: %d, expected: %d", got, expected)
			}

			got = box.textStructure[1]
			expected = len("string")

			if got != expected {
				t.Errorf("Mismatched number of characters in first line: %d, expected: %d", got, expected)
			}
		}
	})

	t.Run("Should parse spaces correctly", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, "a b")

		got := box.RuneAt(1)
		expected := '·'

		if got != expected {
			t.Errorf("Space parsed incorrectly, got: %q, expected %q", got, expected)
		}
	})

	t.Run("Should parse newline correctly", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, "a\nb")

		got := box.RuneAt(1)
		expected := '↩'

		if got != expected {
			t.Errorf("Newline parsed incorrectly, got: %q, expected %q", got, expected)
		}
	})

	t.Run("Should parse tabs correctly", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, "a\tb")

		got := box.RuneAt(1)
		expected := '⇥'

		if got != expected {
			t.Errorf("Newline parsed incorrectly, got: %q, expected %q", got, expected)
		}
	})

	t.Run("Should clear previous text", func(t *testing.T) {
		text := "First string"
		box := NewWpmBox(0, 0, 0, 0, text)
		box.SetText("Second string")

		got := len(box.textStructure)
		expected := 1

		if got != expected {
			t.Errorf("Multiple lines in box: %d, expected: %d", got, expected)
			return
		}
	})

	t.Run("Cursor should be positioned on origin at the start", func(t *testing.T) {
		text := "Example text"
		box := NewWpmBox(0, 0, 0, 0, text)

		if box.line != 0 || box.cursor != 0 {
			t.Errorf("Cursor not in origin position: (%d, %d) expected (0, 0)", box.cursor, box.line)
		}
	})

	t.Run("cursor should be repositioned to start when text is changed", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, "Test")
		box.InsKey('T')
		box.SetText("test2")

		if box.line != 0 || box.cursor != 0 {
			t.Errorf("Cursor not in origin position: (%d, %d) expected (0, 0)", box.cursor, box.line)
		}
	})

	t.Run("Cursor should normally advance horizontally", func(t *testing.T) {
		box := NewWpmBox(0, 0, 0, 0, "Example")
		box.incCursor()

		if box.line != 0 || box.cursor != 1 {
			t.Errorf("Wrong cursor position: (%d, %d) expected (1, 0)", box.cursor, box.line)
		}
	})

	t.Run("incCursor", func(t *testing.T) {
		t.Run("Cursor should follow new line", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 4, "\n")
			box.incCursor()
			if box.line != 1 || box.cursor != 0 {
				t.Errorf("Wrong cursor position: (%d, %d) expected (0, 1)", box.cursor, box.line)
			}
		})

		t.Run("Cursor should follow new line when line is ended", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 4, "test\nline")

			box.InsKey('t')
			box.InsKey('e')
			box.InsKey('s')
			box.InsKey('t')
			box.InsKey('\n')

			expectedX := 0
			expectedY := 1

			if box.line != expectedY || box.cursor != expectedX {
				t.Errorf("Wrong cursor position: (%d, %d) expected (%d, %d)", box.cursor, box.line, expectedX, expectedY)
			}
		})
		t.Run("Should not go to new line before last char", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, "test \nline")
			box.cursor = 4

			box.incCursor()

			expectedX := 5
			expectedY := 0
			if box.line != expectedY || box.cursor != expectedX {
				t.Errorf("Wrong cursor position: (%d, %d) expected (%d, %d)", box.cursor, box.line, expectedX, expectedY)
			}
		})

		t.Run("Should not increment function when text is ended", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, "test")

			box.InsKey('t')
			box.InsKey('e')
			box.InsKey('s')
			box.InsKey('t')
			box.InsKey(' ')

			if box.cursor != 4 || box.line != 0 {
				t.Errorf("Changed cursor position: (%d, %d), expected (4, 0)", box.cursor, box.line)
				return
			}
		})

		t.Run("InsKey behaviour on newline inserted", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, "test\na")
			box.cursor = 3

			box.InsKey('t')

			if box.line != 0 || box.cursor != 4 {
				t.Errorf("Wrong cursor position: (%d, %d) expected (0, 4)", box.line, box.cursor)
			}

		})
	})

	t.Run("decCursor", func(t *testing.T) {
		t.Run("Should normally decremnt cursor position on decCursor called", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, "test")
			box.cursor = 2

			box.decCursor()
			if box.cursor != 1 {
				t.Errorf("Unexpected cursor position: %d, expected 1", box.cursor)
			}
		})

		t.Run("Should not decrement cursor at beginning of text", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 0, "test")

			box.decCursor()
			if box.cursor != 0 {
				t.Errorf("Unexpected cursor position: %d, expected 1", box.cursor)
			}
		})

		t.Run("Should return to previous line on backpress", func(t *testing.T) {
			box := NewWpmBox(0, 0, 0, 4, "test\n")
			box.InsKey('t')
			box.InsKey('e')
			box.InsKey('s')
			box.InsKey('t')
			box.InsKey('\n')

			box.decCursor()
			if box.cursor != 4 || box.line != 0 {
				t.Errorf("Unexpected cursor position: (%d, %d), expected (4, 0)", box.cursor, box.line)
			}
		})
	})

	t.Run("Scrolloff", func(t *testing.T) {
		t.Run("Should increase offset if cursor exceeds scrolloff", func(t *testing.T) {
			// container of height 5 with scrolloff 3
			box := NewWpmBox(0, 0, 0, 5, "\n\n\n\n\n\n\n\n")
			box.ScrollOff = 3

			box.InsKey('\n')
			box.InsKey('\n')
			box.InsKey('\n')

			got := box.offset
			expect := 1
			if got != expect {
				t.Errorf("Wrong offset: %d, expected %d", got, expect)
			}

			got = box.line
			expect = 2

			if got != expect {
				t.Errorf("Unexpected cursor position: %d, expected %d", got, expect)
			}

		})
	})
}
