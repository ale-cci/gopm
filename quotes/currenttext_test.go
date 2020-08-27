package quotes

import "testing"

func TestCurrent(t *testing.T) {
	t.Run("Should correctly intialize", func(t *testing.T) {
		currentText := CurrentText{Text: "Test", Author: "Myself"}

		if currentText.correct != 0 || currentText.wrong != 0 {
			t.Errorf("Wrong CurrentText intialization")
		}
	})

	t.Run("Untyped characters should be blank", func(t *testing.T) {
		currentText := CurrentText{Text: "default"}

		got := currentText.CharStatus(0)
		expected := BLANK

		if got != expected {
			t.Errorf("Wrong character staus for unvisited character: %v, expected %v", got, expected)
		}
	})

	t.Run("Right cells should be right", func(t *testing.T) {
		ct := CurrentText{Text: "test"}
		ct.correct = 2

		got := ct.CharStatus(1)
		expected := RIGHT

		if got != expected {
			t.Errorf("Wrong character staus for right character: %v, expected %v", got, expected)
		}
	})

	t.Run("Wrong cells should be wrong", func(t *testing.T) {
		ct := CurrentText{Text: "test"}
		ct.wrong = 2

		got := ct.CharStatus(1)
		expected := WRONG

		if got != expected {
			t.Errorf("Wrong character staus for wrong character: %v, expected %v", got, expected)
		}
	})

	t.Run("Test mixed cell status", func(t *testing.T) {
		ct := CurrentText{Text: "RWB"}
		ct.correct = 1
		ct.wrong = 1

		got := ct.CharStatus(0)
		expected := RIGHT

		if got != expected {
			t.Errorf("Wrong character staus for character: %v, expected %v", got, expected)
		}

		got = ct.CharStatus(1)
		expected = WRONG

		if got != expected {
			t.Errorf("Wrong character staus for character: %v, expected %v", got, expected)
		}

		got = ct.CharStatus(2)
		expected = BLANK

		if got != expected {
			t.Errorf("Wrong character staus for character: %v, expected %v", got, expected)
		}
	})

	t.Run("InsKey", func(t *testing.T) {
		t.Run("Should insert detect correct key when inserted", func(t *testing.T) {
			ct := CurrentText{Text: "Abcd"}
			ct.InsKey('A')

			got := ct.correct
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should insert detect correct key when inserted", func(t *testing.T) {
			ct := CurrentText{Text: "Abcd"}
			ct.InsKey('B')

			got := ct.wrong
			expected := 1

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Detect correct characters in succession", func(t *testing.T) {
			ct := CurrentText{Text: "Test"}
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
			ct := CurrentText{Text: "s"}
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
			ct := CurrentText{Text: "s"}
			ct.InsKey('s')
			ct.Backspace()

			got := ct.correct
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should not delete characters at start of function", func(t *testing.T) {
			ct := CurrentText{Text: "s"}
			ct.Backspace()

			got := ct.correct
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of correct characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should delete wrong characters", func(t *testing.T) {
			ct := CurrentText{Text: "s"}
			ct.InsKey('a')
			ct.Backspace()

			got := ct.wrong
			expected := 0

			if got != expected {
				t.Errorf("Wrong number of wrong characters: %d expected: %d", got, expected)
			}
		})

		t.Run("Should delete wrong characters first", func(t *testing.T) {
			ct := CurrentText{Text: "Test string"}
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
}
