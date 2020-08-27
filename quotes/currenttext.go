package quotes

type CurrentText struct {
	Text   string
	Author string

	// Correct and wrong characters typed
	correct, wrong int
}

func (t *CurrentText) CurrentRune() rune {
	return t.RuneAt(t.Position())
}

func (t *CurrentText) RuneAt(position int) rune {
	return []rune(t.Text)[position]
}

func (t *CurrentText) InsKey(char rune) {
	if !t.IsEndPosition() {
		current := t.CurrentRune()

		if char == current && t.wrong == 0 {
			t.correct += 1
		} else {
			t.wrong += 1
		}
	}
}

func (t *CurrentText) Backspace() {
	if !t.IsStartPosition() {
		if t.wrong > 0 {
			t.wrong -= 1
		} else {
			t.correct -= 1
		}
	}
}

func (t *CurrentText) CharStatus(position int) CellState {
	switch {
	case position < t.correct:
		return RIGHT
	case position >= t.correct && position < t.correct+t.wrong:
		return WRONG
	default:
		return BLANK
	}
}

func (t *CurrentText) Position() int {
	return t.correct + t.wrong
}

func (t *CurrentText) IsStartPosition() bool {
	return t.Position() == 0
}

func (t *CurrentText) IsEndPosition() bool {
	return t.Position() == len(t.Text)
}
