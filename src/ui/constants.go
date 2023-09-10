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
)

type BlockOperation int

const (
	SELECT BlockOperation = iota
	PLACE_BLUE
	PLACE_RED
	PLACE_YELLOW
)
