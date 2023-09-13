package ui

type Renderer int

const (
	ISOMETRIC Renderer = iota
	TWO_DIMENSIONAL
)

type BlockSize int

const (
	HALF BlockSize = iota
	FULL
	FLAT
)

func (blockSize BlockSize) GetHeight() int {
	switch blockSize {
	case FULL:
		return 2
	case HALF:
		return 1
	default:
		return 0
	}
}

type BlockOperation int

const (
	SELECT BlockOperation = iota
	PLACE_BLUE
	PLACE_RED
	PLACE_YELLOW
)
