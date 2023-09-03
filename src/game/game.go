package game

import (
	"image/color"

	"github.com/timothy-ch-cheung/go-game-block-placement/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Renderer int

const (
	ISOMETRIC Renderer = iota
	TWO_DIMENSIONAL
)

type Game struct {
	loader        *resource.Loader
	background    *ebiten.Image
	renderingMode Renderer
}

func NewGame() *Game {
	g := &Game{}

	audioContext := audio.NewContext(44100)
	loader := resource.NewLoader(audioContext)
	loader.OpenAssetFunc = assets.OpenAssetFunc
	assets.RegisterImageResources(loader)
	g.loader = loader

	background := ebiten.NewImage(ScreenWidth, ScreenHeight)
	background.Fill(color.RGBA{R: 21, G: 29, B: 40, A: 1}) // #151d28
	g.background = background

	g.renderingMode = ISOMETRIC

	return g
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})

	switch g.renderingMode {
	case ISOMETRIC:
		renderIso(g, screen)
	case TWO_DIMENSIONAL:
		render2D(g, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
