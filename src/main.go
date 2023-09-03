package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/timothy-ch-cheung/go-game-block-placement/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*game.Scale, game.ScreenHeight*game.Scale)
	ebiten.SetWindowTitle("Game Block Placement Demo")

	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
