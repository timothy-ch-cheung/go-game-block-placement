package game

import (
	"image/color"

	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/objects"
	"github.com/timothy-ch-cheung/go-game-block-placement/ui"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"

	input "github.com/quasilyte/ebitengine-input"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game struct {
	inputSystem  input.System
	inputHandler *input.Handler
	loader       *resource.Loader
	background   *ebiten.Image
	board        *objects.Board
	ui           *ui.UI
	cursor       *resolv.Object
}

func NewGame() *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	g.inputHandler = g.inputSystem.NewHandler(0, ui.NewKeyMap())

	audioContext := audio.NewContext(44100)
	loader := resource.NewLoader(audioContext)
	loader.OpenAssetFunc = assets.OpenAssetFunc
	assets.RegisterImageResources(loader)
	g.loader = loader

	background := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	background.Fill(color.RGBA{R: 21, G: 29, B: 40, A: 1}) // #151d28
	g.background = background

	x, y := ebiten.CursorPosition()
	g.cursor = resolv.NewObject(float64(x), float64(y), 1, 1)
	g.board = objects.NewBoard(10, 10, 5, g.cursor, loader)

	var viewModeChangedHandler widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if g.ui.State.Renderer == ui.ISOMETRIC {
			g.ui.State.Renderer = ui.TWO_DIMENSIONAL
		} else {
			g.ui.State.Renderer = ui.ISOMETRIC
		}
	}
	var blockSizeChangedHandler widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if g.ui.State.BlockSize == ui.HALF {
			g.ui.State.BlockSize = ui.FULL
		} else {
			g.ui.State.BlockSize = ui.HALF
		}
	}

	handlers := &ui.Handlers{
		ViewToggleChangedHandler: &viewModeChangedHandler,
		BlockSizeChangedHandler:  &blockSizeChangedHandler,
	}

	g.ui = ui.NewUserInterface(handlers, loader)

	return g
}

func (g *Game) Update() error {
	g.ui.Update()
	g.inputSystem.Update()
	g.board.Update(g.ui.State, g.inputHandler)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
	g.ui.Draw(screen)
	switch g.ui.State.Renderer {
	case ui.ISOMETRIC:
		g.board.RenderIso(screen)
	case ui.TWO_DIMENSIONAL:
		g.board.Render2D(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
