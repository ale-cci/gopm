package tui

import (
	"strings"
)

type ParsedText struct {
	Structure []int
	Text      []rune
}

func (pt *ParsedText) RuneAt(position int) rune {
	return pt.Text[position]
}

func NewParsedText(Text string) *ParsedText {
	pt := ParsedText{Text: []rune(Text)}

	lines := strings.Split(Text, "\n")
	pt.Structure = make([]int, len(lines))

	for i, line := range lines {
		for _, char := range line {
			pt.Structure[i] += runeSize(char)
		}

		// new line character to all lines except last one
		if i != len(lines)-1 {
			pt.Structure[i] += 1
		}
	}

	return &pt
}
