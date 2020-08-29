package wpm

import (
	"math"
	"testing"
	"time"
)

func TestCurrent(t *testing.T) {
	t.Run("Should correctly initialize ", func(t *testing.T) {
		currentText := KeystrokeCounter{Text: "Test"}

		if currentText.correct != 0 || currentText.wrong != 0 {
			t.Errorf("Wrong KeystrokeCounter initialize")
		}
	})

	t.Run("Untyped characters should be blank", func(t *testing.T) {
		currentText := KeystrokeCounter{Text: "default"}

		got := currentText.CharStatus(0)
		expected := BLANK

		if got != expected {
			t.Errorf("Wrong character status for unvisited character: %v, expected %v", got, expected)
		}
	})

	t.Run("Right cells should be right", func(t *testing.T) {
		ct := KeystrokeCounter{Text: "test"}
		ct.correct = 2

		got := ct.CharStatus(1)
		expected := RIGHT

		if got != expected {
			t.Errorf("Wrong character status for right character: %v, expected %v", got, expected)
		}
	})

	t.Run("Wrong cells should be wrong", func(t *testing.T) {
		ct := KeystrokeCounter{Text: "test"}
		ct.wrong = 2

		got := ct.CharStatus(1)
		expected := WRONG

		if got != expected {
			t.Errorf("Wrong character status for wrong character: %v, expected %v", got, expected)
		}
	})

	t.Run("Test mixed cell status", func(t *testing.T) {
		ct := KeystrokeCounter{Text: "RWB"}
		ct.correct = 1
		ct.wrong = 1

		got := ct.CharStatus(0)
		expected := RIGHT

		if got != expected {
			t.Errorf("Wrong character status for character: %v, expected %v", got, expected)
		}

		got = ct.CharStatus(1)
		expected = WRONG

		if got != expected {
			t.Errorf("Wrong character status for character: %v, expected %v", got, expected)
		}

		got = ct.CharStatus(2)
		expected = BLANK

		if got != expected {
			t.Errorf("Wrong character status for character: %v, expected %v", got, expected)
		}
	})

	t.Run("InsKey", func(t *testing.T) {
		t.Run("InsKey should identify correct key", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "Abcd"}
			ct.InsKey('A')

			got := ct.correct
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should insert detect correct key when inserted", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "Abcd"}
			ct.InsKey('B')

			got := ct.wrong
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}
		})

		t.Run("On error all next keystrokes should be error", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "Abcd"}
			ct.InsKey('B')
			ct.InsKey('b')
			ct.InsKey('c')

			got := ct.wrong
			expected := 3

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Detect correct characters in succession", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "Test"}
			ct.InsKey('T')
			ct.InsKey('e')
			ct.InsKey('s')

			got := ct.correct
			expected := 3

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should stop character insertion when string is finished", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "s"}
			ct.InsKey('s')
			ct.InsKey('s')

			got := ct.correct
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})
	})

	t.Run("Backspace", func(t *testing.T) {
		t.Run("Should delete correct character inserted on backspace", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "s"}
			ct.InsKey('s')
			ct.Backspace()

			got := ct.correct
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should not delete characters at start of function", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "s"}
			ct.Backspace()

			got := ct.correct
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should delete wrong characters", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "s"}
			ct.InsKey('a')
			ct.Backspace()

			got := ct.wrong
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should delete wrong characters first", func(t *testing.T) {
			ct := KeystrokeCounter{Text: "Test string"}
			ct.InsKey('T')
			ct.InsKey('e')
			ct.InsKey('s')
			ct.InsKey('s')
			ct.InsKey('s')
			ct.Backspace()

			got := ct.wrong
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}

			got = ct.correct
			expected = 3

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})
	})

	t.Run("Wpm", func(t *testing.T) {
		tt := []struct {
			name     string
			duration string
			text     string
			typed    []rune
			expect   float64
		}{
			{
				name:     "60cpm",
				duration: "1s",
				text:     "test",
				typed:    []rune("t"),
				expect:   60,
			},
			{
				name:     "1cpm",
				duration: "1m",
				text:     "a",
				typed:    []rune("abcdefgehij"),
				expect:   1,
			},
			{
				name:     "Start wpm",
				duration: "0s",
				text:     "test",
				typed:    []rune("test"),
				expect:   0,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				start := time.Time{}
				duration, err := time.ParseDuration(tc.duration)

				if err != nil {
					t.Fatalf("Failed to parse duration: %v", err)
				}
				end := start.Add(duration)

				counter := KeystrokeCounter{Text: tc.text, Start: start}
				for _, k := range tc.typed {
					counter.InsKey(k)
				}

				got := counter.Cpm(end)
				expect := tc.expect

				if math.Abs(got-expect) >= 0.001 {
					t.Fatalf("Wrong cpm calculation for %s: %v expected %v", tc.name, got, expect)
				}
			})
		}
	})
}
