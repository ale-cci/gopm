package quotes

// possible cell status
type CharStatus int

const (
	WRONG CharStatus = iota
	RIGHT CharStatus = iota
	BLANK CharStatus = iota
)

func (c CharStatus) String() string {
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
