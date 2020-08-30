package tui

import (
	"testing"

	"github.com/ale-cci/gopm/wpm"
)

func TestWpmBox(t *testing.T) {

	/*
		t.Run("Should parse spaces correctly", func(t *testing.T) {
			got := parseRune(' ')
			expected := '·'

			if got != expected {
				t.Errorf("Space parsed incorrectly, got: %q, expected %q", got, expected)
			}
		})

		t.Run("Should parse newline correctly", func(t *testing.T) {
			got := parseRune('\n')
			expected := '↩'

			if got != expected {
				t.Errorf("Newline parsed incorrectly, got: %q, expected %q", got, expected)
			}
		})

		t.Run("Should parse tabs correctly", func(t *testing.T) {
			got := parseRune('\t')
			expected := '⇥'

			if got != expected {
				t.Errorf("Newline parsed incorrectly, got: %q, expected %q", got, expected)
			}
		})
	*/

	t.Run("Should clear previous Text", func(t *testing.T) {
		Text := NewParsedText("First string")
		counter := wpm.KeystrokeCounter{Capacity: 8}
		box := NewWpmBox(0, 0, 0, 0, Text, &counter)

		box.IncCursor()
		counter.InsKey(true)
		box.IncCursor()
		counter.InsKey(true)
		box.IncCursor()
		counter.InsKey(true)

		box.SetText(NewParsedText("Second"))

		tt := []struct {
			name string
			val  int
		}{
			{"cursor", box.cursor},
			{"line", box.line},
			{"offset", box.offset},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.val
				expected := 0

				if got != expected {
					t.Errorf("%s has not been resetted, got: %d, expected: %d", tc.name, got, expected)
					return
				}
			})
		}
	})

	t.Run("IncCursor", func(t *testing.T) {
		tt := []struct {
			name      string
			text      string
			height    int
			incr      int
			decr      int
			cursor    int
			line      int
			offset    int
			scrolloff int
		}{
			{
				name:   "Cursor should normally advance horizontally",
				text:   "Example",
				height: 4,
				incr:   1,
				cursor: 1,
				line:   0,
			},
			{
				name:   "Cursor should follow new line",
				text:   "\n",
				height: 4,
				incr:   1,
				cursor: 0,
				line:   1,
			},
			{
				name:   "cursor should follow new line when line is ended",
				text:   "test\nline",
				height: 4,
				incr:   5,
				cursor: 0,
				line:   1,
			},
			{
				name:   "should not go to new line before last char",
				text:   "test \nline",
				height: 4,
				incr:   5,
				cursor: 5,
				line:   0,
			},
			{
				name:   "Should not increment when text is ended",
				text:   "test",
				height: 3,
				incr:   5,
				cursor: 4,
				line:   0,
			},
			{
				name:   "Should advance cursor more if tab inserted",
				text:   "A\tB",
				height: 4,
				incr:   2,
				cursor: 5,
				line:   0,
			},
			{
				name:   "Should normally decrement cursor position on DecCursor called",
				text:   "test",
				height: 1,
				incr:   2,
				decr:   1,
				cursor: 1,
				line:   0,
			},
			{
				name:   "should not decrement cursor at beginning of text",
				text:   "test",
				height: 1,
				incr:   0,
				decr:   1,
				cursor: 0,
				line:   0,
			},
			{
				name:   "Should return to previous line on backpress",
				text:   "test\n",
				height: 4,
				incr:   5,
				decr:   1,
				cursor: 4,
				line:   0,
			},
			{
				name:   "Should return to original position after tab delete",
				text:   "a\t",
				height: 4,
				incr:   2,
				decr:   2,
				cursor: 0,
				line:   0,
			},
			{
				name:   "Should set cursor to line end when previous line has characters of length > 1",
				text:   "b\t\n",
				height: 6,
				incr:   3,
				decr:   1,
				cursor: 5,
				line:   0,
			},
			{
				name:      "no offset increase at end of page",
				text:      "a\n\n",
				scrolloff: 3,
				height:    20,
				incr:      2,
				decr:      1,
				offset:    0,
				cursor:    1,
				line:      0,
			},
			{
				name:      "offset increase when text out of screen",
				text:      "\n\n\n\n\n\n\n\n",
				height:    5,
				scrolloff: 3,
				incr:      3,
				offset:    1,
				cursor:    0,
				line:      2,
			},
			{
				name:      "remove offset when cursor moves up",
				text:      "\n\n\n\n\n\n\n",
				height:    6,
				scrolloff: 3,
				incr:      4,
				decr:      1,
				offset:    0,
				cursor:    0,
				line:      3,
			},
			{
				name:      "scroll until end of text",
				text:      "\n\n\n\n\n\n\n\n",
				height:    6,
				scrolloff: 4,
				incr:      8,
				decr:      0,
				offset:    2,
				cursor:    0,
				line:      6,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				counter := wpm.KeystrokeCounter{Capacity: len(tc.text)}
				text := NewParsedText(tc.text)
				box := NewWpmBox(0, 0, 0, tc.height, text, &counter)
				box.ScrollOff = tc.scrolloff

				for i := 0; i < tc.incr; i++ {
					box.IncCursor()
					counter.InsKey(true)
				}

				for i := 0; i < tc.decr; i++ {
					box.DecCursor()
					counter.Backspace()
				}

				if box.line != tc.line || box.cursor != tc.cursor {
					t.Errorf("%d, Wrong cursor position: (%d, %d) expected (%d, %d)", counter.Inserted(), box.cursor, box.line, tc.cursor, tc.line)
				}

				if box.offset != tc.offset {
					t.Fatalf("Wrong offset for %s: %d, expected: %d", tc.name, box.offset, tc.offset)
				}
			})
		}
	})

	/*
			t.Run("Should not set negative scrolloff", func(t *testing.T) {
				box := NewWpmBox(0, 0, 0, 20, "a\n\n")
				box.ScrollOff = 3
				box.IncCursor()
				box.IncCursor()
				box.DecCursor()

				expect := box.offset
				got := 0

				if got != expect {
					t.Errorf("Wrong offset value: %d, expected %d", got, expect)
				}
			})
		t.Run("Scrolloff", func(t *testing.T) {
			t.Run("Should increase offset if cursor exceeds scrolloff", func(t *testing.T) {
				// container of height 5 with scrolloff 3
				box := NewWpmBox(0, 0, 0, 5, "\n\n\n\n\n\n\n\n")
				box.ScrollOff = 3

				box.IncCursor()
				box.IncCursor()
				box.IncCursor()

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

			t.Run("Should decrease offset correctly if cursor moves up", func(t *testing.T) {
				box := NewWpmBox(0, 0, 0, 6, "\n\n\n\n\n\n\n")
				box.ScrollOff = 3
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.DecCursor()

				got := box.offset
				expect := 0

				if got != expect {
					t.Errorf("Wrong offset: %d, expected: %d", got, expect)
				}
			})

			t.Run("Should scroll till the end of Text", func(t *testing.T) {
				box := NewWpmBox(0, 0, 0, 6, "\n\n\n\n\n\n\n\n")
				box.ScrollOff = 4
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()
				box.IncCursor()

				got := box.offset
				expect := 2

				if got != expect {
					t.Errorf("Wrong offset: %d, expected: %d", got, expect)
				}
			})
		})
	*/
}
