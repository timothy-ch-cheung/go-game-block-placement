package ui

import (
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionSelect input.Action = iota
	ActionDelete
)

func NewKeyMap() input.Keymap {
	return input.Keymap{
		ActionSelect: {input.KeyMouseLeft},
		ActionDelete: {input.KeyMouseRight},
	}
}
