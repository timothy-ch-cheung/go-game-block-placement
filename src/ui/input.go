package ui

import (
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionSelect input.Action = iota
)

func NewKeyMap() input.Keymap {
	return input.Keymap{
		ActionSelect: {input.KeyMouseLeft},
	}
}
