package game

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/objects"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"

	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Renderer int

const (
	ISOMETRIC Renderer = iota
	TWO_DIMENSIONAL
)

type Game struct {
	inputSystem   *input.System
	loader        *resource.Loader
	background    *ebiten.Image
	renderingMode Renderer
	board         *objects.Board
	ui            *ebitenui.UI
	cursor        *resolv.Object
}

func NewGame() *Game {
	g := &Game{}

	g.inputSystem.NewHandler(0, NewKeyMap())

	audioContext := audio.NewContext(44100)
	loader := resource.NewLoader(audioContext)
	loader.OpenAssetFunc = assets.OpenAssetFunc
	assets.RegisterImageResources(loader)
	g.loader = loader

	background := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	background.Fill(color.RGBA{R: 21, G: 29, B: 40, A: 1}) // #151d28
	g.background = background

	g.renderingMode = ISOMETRIC

	x, y := ebiten.CursorPosition()
	g.cursor = resolv.NewObject(float64(x), float64(y), 1, 1)
	g.board = objects.NewBoard(10, 10, 5, g.cursor, loader)

	var viewModeChangedHandler widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if g.renderingMode == ISOMETRIC {
			g.renderingMode = TWO_DIMENSIONAL
			args.Active.SetState(1)
		} else {
			g.renderingMode = ISOMETRIC
			args.Active.SetState(0)
		}
	}
	handlers := &Handlers{
		viewToggleChangedHandler: &viewModeChangedHandler,
	}

	g.ui = newUserInterface(handlers, loader)

	return g
}

func (g *Game) Update() error {
	g.ui.Update()
	g.inputSystem.Update()
	x, y := ebiten.CursorPosition()
	g.cursor.X = float64(x)
	g.cursor.Y = float64(y)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
	g.ui.Draw(screen)
	switch g.renderingMode {
	case ISOMETRIC:
		g.board.RenderIso(screen)
	case TWO_DIMENSIONAL:
		g.board.Render2D(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
