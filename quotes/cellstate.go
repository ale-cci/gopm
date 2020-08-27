package quotes

// possible cell status
type CellState int

const (
	WRONG CellState = iota
	RIGHT CellState = iota
	BLANK CellState = iota
)

func (c CellState) String() string {
	switch c {
	case WRONG:
		return "WRONG"
	case RIGHT:
		return "RIGHT"
	case BLANK:
		return "BLANK"
	default:
		return "???"
	}
}
