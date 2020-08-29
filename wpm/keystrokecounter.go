package wpm

import "time"

type KeystrokeCounter struct {
	Text            string
	Right, Wrong    int
	Start           time.Time
	TotalWrong      int
	TotalKeystrokes int
}

func (t *KeystrokeCounter) CurrentRune() rune {
	return t.RuneAt(t.Position())
}

func (t *KeystrokeCounter) RuneAt(position int) rune {
	return []rune(t.Text)[position]
}

func (t *KeystrokeCounter) InsKey(char rune) {

	if !t.IsEndPosition() {
		current := t.CurrentRune()

		if char == current && t.Wrong == 0 {
			t.Right += 1
		} else {
			t.Wrong += 1
			t.TotalWrong += 1
		}
		t.TotalKeystrokes += 1
	}
}

func (t *KeystrokeCounter) Cpm(now time.Time) float64 {
	elapsed := now.Sub(t.Start)
	typed := t.Right
	time := elapsed.Minutes()

	if time < 0.00001 {
		return 0
	}
	return float64(typed) / time
}

func (t *KeystrokeCounter) Wpm(now time.Time) float64 {
	return t.Cpm(now) / 5
}

func (t *KeystrokeCounter) Accuracy() float64 {
	return float64(t.TotalKeystrokes-t.TotalWrong) * 100 / float64(t.TotalKeystrokes)
}

// Remove the last character inserted
func (t *KeystrokeCounter) Backspace() {
	if !t.IsStartPosition() {
		if t.Wrong > 0 {
			t.Wrong -= 1
		} else {
			t.Right -= 1
		}
	}
}

func (t *KeystrokeCounter) Position() int {
	return t.Right + t.Wrong
}

func (t *KeystrokeCounter) IsStartPosition() bool {
	return t.Position() == 0
}

// Has reached the end of the text
func (t *KeystrokeCounter) IsEndPosition() bool {
	return t.Position() == len(t.Text)
}

// Return the status of the character at `position`
// RIGHT, WRONG or BLANK if the character has not been typed yet
func (t *KeystrokeCounter) CharStatus(position int) CharStatus {
	switch {
	case position < t.Right:
		return RIGHT
	case position >= t.Right && position < t.Right+t.Wrong:
		return WRONG
	default:
		return BLANK
	}
}
