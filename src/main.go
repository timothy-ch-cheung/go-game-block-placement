package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	resource "github.com/quasilyte/ebitengine-resource"
)

type game struct {
	loader *resource.Loader
	background *ebiten.Image
}

func newGame() *game {
	g := &game{}

	background := ebiten.NewImage(ScreenWidth, ScreenHeight)
	background.Fill(color.RGBA{R: 21, G: 29, B: 40, A: 1}) // #151d28
	g.background = background

	return g
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*Scale, ScreenHeight*Scale)
	ebiten.SetWindowTitle("Game Block Placement Demo")

	game := newGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
