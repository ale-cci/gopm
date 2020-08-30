package wpm

import (
	"math"
	"testing"
	"time"
)

func TestCurrent(t *testing.T) {
	t.Run("CharStatus", func(t *testing.T) {
		tt := []struct {
			name       string
			keystrokes []bool
			position   int
			status     CharStatus
		}{
			{
				name:       "untyped character",
				keystrokes: []bool{},
				position:   3,
				status:     BLANK,
			},
			{
				name:       "right character",
				keystrokes: []bool{true},
				position:   0,
				status:     RIGHT,
			},
			{
				name:       "wrong character",
				keystrokes: []bool{false},
				position:   0,
				status:     WRONG,
			},
			{
				name:       "wrong after right character",
				keystrokes: []bool{true, false},
				position:   1,
				status:     WRONG,
			},
			{
				name:       "right after wrong",
				keystrokes: []bool{false, true},
				position:   1,
				status:     WRONG,
			},
			{
				name:       "right before wrong",
				keystrokes: []bool{true, false},
				position:   0,
				status:     RIGHT,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				counter := KeystrokeCounter{Capacity: 8}

				for _, k := range tc.keystrokes {
					counter.InsKey(k)
				}

				got := counter.CharStatus(tc.position)
				expect := tc.status

				if got != expect {
					t.Errorf("Unexpected status for %s: %s, expected: %s", tc.name, got, expect)
				}
			})
		}
	})

	t.Run("InsKey and Backspace", func(t *testing.T) {
		type Action int
		const (
			A_RIGHT     Action = iota
			A_WRONG     Action = iota
			A_BACKSPACE Action = iota
		)

		tt := []struct {
			name       string
			keystrokes []Action
			wrong      int
			right      int
			totalWrong int
			capacity   int
		}{
			{
				name:       "Should stop character insertion when string is finished",
				keystrokes: []Action{A_RIGHT, A_WRONG},
				wrong:      0,
				right:      1,
				totalWrong: 0,
				capacity:   1,
			},
			{
				name:       "Detect correct characters in succession",
				keystrokes: []Action{A_RIGHT, A_RIGHT, A_RIGHT},
				wrong:      0,
				right:      3,
				totalWrong: 0,
				capacity:   4,
			},
			{
				name:       "On error all next keystrokes should be errors",
				capacity:   4,
				keystrokes: []Action{A_WRONG, A_RIGHT, A_RIGHT},
				wrong:      3,
				right:      0,
				totalWrong: 3,
			},
			{
				name:       "Should delete wrong characters first",
				keystrokes: []Action{A_RIGHT, A_RIGHT, A_RIGHT, A_WRONG, A_WRONG, A_BACKSPACE},
				wrong:      1,
				right:      3,
				totalWrong: 2,
				capacity:   8,
			},
			{
				name:       "Should delete Wrong characters",
				keystrokes: []Action{A_WRONG, A_BACKSPACE},
				wrong:      0,
				right:      0,
				totalWrong: 1,
				capacity:   1,
			},
			{
				name:       "should delete all right and wrong characters",
				keystrokes: []Action{A_RIGHT, A_RIGHT, A_WRONG, A_RIGHT, A_WRONG, A_BACKSPACE, A_BACKSPACE, A_BACKSPACE, A_BACKSPACE, A_BACKSPACE},
				wrong:      0,
				right:      0,
				totalWrong: 3,
				capacity:   8,
			},
			{
				name:       "Should not delete characters at start of function",
				keystrokes: []Action{A_BACKSPACE},
				wrong:      0,
				right:      0,
				totalWrong: 0,
				capacity:   1,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ct := KeystrokeCounter{Capacity: tc.capacity}

				for _, action := range tc.keystrokes {
					switch action {
					case A_RIGHT:
						ct.InsKey(true)
					case A_WRONG:
						ct.InsKey(false)
					case A_BACKSPACE:
						ct.Backspace()
					}
				}

				if tc.wrong != ct.Wrong {
					t.Errorf("%s, wrong characters: %d, expected %d", tc.name, ct.Wrong, tc.wrong)
				}

				if tc.right != ct.Right {
					t.Errorf("%s, right characters: %d, expected %d", tc.name, ct.Right, tc.right)
				}

				if tc.totalWrong != ct.TotalWrong {
					t.Errorf("%s, total wrong characters: %v, expected %v", tc.name, ct.TotalWrong, tc.totalWrong)
				}
			})
		}
	})

	t.Run("Wpm", func(t *testing.T) {
		tt := []struct {
			name     string
			duration string
			correct  int
			expect   float64
		}{
			{
				name:     "60cpm",
				duration: "1s",
				correct:  1,
				expect:   60,
			},
			{
				name:     "1cpm",
				duration: "1m",
				correct:  1,
				expect:   1,
			},
			{
				name:     "Start wpm",
				duration: "0s",
				correct:  4,
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

				counter := KeystrokeCounter{Right: tc.correct, Start: start}

				got := counter.Cpm(end)
				expect := tc.expect

				if math.Abs(got-expect) >= 0.001 {
					t.Fatalf("Wrong cpm calculation for %s: %v expected %v", tc.name, got, expect)
				}
			})
		}
	})
}
