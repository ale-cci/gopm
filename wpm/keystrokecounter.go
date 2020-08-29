package wpm

import "time"

type KeystrokeCounter struct {
	Right, Wrong    int
	Start           time.Time
	TotalWrong      int
	TotalKeystrokes int
	Capacity        int
}

func (t *KeystrokeCounter) InsKey(correct bool) {

	if !t.IsEndPosition() {
		if correct && t.Wrong == 0 {
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

func (t *KeystrokeCounter) IsEndPosition() bool {
	return t.Position() == t.Capacity
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

func (t *KeystrokeCounter) Reset() {
	t.Right = 0
	t.Wrong = 0
	t.TotalWrong = 0
	t.TotalKeystrokes = 0
}
