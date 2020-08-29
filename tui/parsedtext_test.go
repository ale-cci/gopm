package tui_test

import (
	"testing"

	tui "github.com/ale-cci/gopm/tui"
)

func TestParsedText(t *testing.T) {
	t.Run("Should correctly initialize buffer", func(t *testing.T) {
		pt := tui.NewParsedText("multiline\nstring")

		got := len(pt.Structure)
		expected := 2

		if got != expected {
			t.Errorf("Different line of Text parsed: %d, expected: %d", got, expected)
		} else {

			got = pt.Structure[0]
			expected = len("multiline\n")

			if got != expected {
				t.Errorf("Mismatched number of characters in first line: %d, expected: %d", got, expected)
			}

			got = pt.Structure[1]
			expected = len("string")

			if got != expected {
				t.Errorf("Mismatched number of characters in first line: %d, expected: %d", got, expected)
			}
		}
	})
}
