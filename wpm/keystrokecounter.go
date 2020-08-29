package wpm

import "time"

type KeystrokeCounter struct {
	Text string
	// Correct and wrong characters typed
	correct, wrong int
	Start          time.Time
}

func (t *KeystrokeCounter) CurrentRune() rune {
	return t.RuneAt(t.Position())
}

func (t *KeystrokeCounter) RuneAt(position int) rune {
	return []rune(t.Text)[position]
}

func (t *KeystrokeCounter) InsKey(char rune) {
	if t.IsStartPosition() {
		t.Start = time.Now()
	}

	if !t.IsEndPosition() {
		current := t.CurrentRune()

		if char == current && t.wrong == 0 {
			t.correct += 1
		} else {
			t.wrong += 1
		}
	}
}

func (t *KeystrokeCounter) Cpm(now time.Time) float64 {
	elapsed := now.Sub(t.Start)
	typed := t.correct
	time := elapsed.Minutes()

	if time < 0.00001 {
		return 0
	}
	return float64(typed) / time
}

func (t *KeystrokeCounter) Wpm(now time.Time) float64 {
	return t.Cpm(now) / 5
}

// Remove the last character inserted
func (t *KeystrokeCounter) Backspace() {
	if !t.IsStartPosition() {
		if t.wrong > 0 {
			t.wrong -= 1
		} else {
			t.correct -= 1
		}
	}
}

// Return the status of the character at `position`
// RIGHT, WRONG or BLANK if the character has not been typed yet
func (t *KeystrokeCounter) CharStatus(position int) CharStatus {
	switch {
	case position < t.correct:
		return RIGHT
	case position >= t.correct && position < t.correct+t.wrong:
		return WRONG
	default:
		return BLANK
	}
}

func (t *KeystrokeCounter) Position() int {
	return t.correct + t.wrong
}

func (t *KeystrokeCounter) IsStartPosition() bool {
	return t.Position() == 0
}

// Has reached the end of the text
func (t *KeystrokeCounter) IsEndPosition() bool {
	return t.Position() == len(t.Text)
}
