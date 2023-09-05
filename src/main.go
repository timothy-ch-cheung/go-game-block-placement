package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/timothy-ch-cheung/go-game-block-placement/game"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth*config.Scale, config.ScreenHeight*config.Scale)
	ebiten.SetWindowTitle("Game Block Placement Demo")

	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
