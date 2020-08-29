package wpm

import (
	"math"
	"testing"
	"time"
)

func TestCurrent(t *testing.T) {
	t.Run("Should Rightly initialize ", func(t *testing.T) {
		currentText := KeystrokeCounter{Text: "Test"}

		if currentText.Right != 0 || currentText.Wrong != 0 {
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
		ct.Right = 2

		got := ct.CharStatus(1)
		expected := RIGHT

		if got != expected {
			t.Errorf("Wrong character status for right character: %v, expected %v", got, expected)
		}
	})

	t.Run("Wrong cells should be Wrong", func(t *testing.T) {
		ct := KeystrokeCounter{Text: "test"}
		ct.Wrong = 2

		got := ct.CharStatus(1)
		expected := WRONG

		if got != expected {
			t.Errorf("Wrong character status for Wrong character: %v, expected %v", got, expected)
		}
	})

	t.Run("Test mixed cell status", func(t *testing.T) {
		ct := KeystrokeCounter{Text: "RWB"}
		ct.Right = 1
		ct.Wrong = 1

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

	t.Run("InsKey and Backspace", func(t *testing.T) {
		tt := []struct {
			name       string
			text       string
			keystrokes []rune
			wrong      int
			right      int
			totalWrong int
		}{
			{
				name:       "Should stop character insertion when string is finished",
				text:       "s",
				keystrokes: []rune("ss"),
				wrong:      0,
				right:      1,
				totalWrong: 0,
			},
			{
				name:       "Detect correct characters in succession",
				text:       "Test",
				keystrokes: []rune("Tes"),
				wrong:      0,
				right:      3,
				totalWrong: 0,
			},
			{
				name:       "On error all next keystrokes should be errors",
				text:       "Abcd",
				keystrokes: []rune("Bbc"),
				wrong:      3,
				right:      0,
				totalWrong: 3,
			},
			{
				name:       "Should delete wrong characters first",
				text:       "Test string",
				keystrokes: []rune("Tesss\b"),
				wrong:      1,
				right:      3,
				totalWrong: 2,
			},
			{
				name:       "Should delete Wrong characters",
				text:       "s",
				keystrokes: []rune("a\b"),
				wrong:      0,
				right:      0,
				totalWrong: 1,
			},
			{
				name:       "Should not delete characters at start of function",
				text:       "s",
				keystrokes: []rune{'\b'},
				wrong:      0,
				right:      0,
				totalWrong: 0,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ct := KeystrokeCounter{Text: tc.text}

				for _, k := range tc.keystrokes {
					if k == '\b' {
						ct.Backspace()
					} else {
						ct.InsKey(k)
					}
				}

				if tc.wrong != ct.Wrong {
					t.Errorf("Mismatched number of wrong characters: %d, expected %d", ct.Wrong, tc.wrong)
				}

				if tc.right != ct.Right {
					t.Errorf("Mismatched number of right characters: %d, expected %d", ct.Right, tc.right)
				}

				if tc.totalWrong != ct.TotalWrong {
					t.Errorf("%q, total wrong characters: %v, expected %v", tc.name, ct.TotalWrong, tc.totalWrong)
				}
			})
		}
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
