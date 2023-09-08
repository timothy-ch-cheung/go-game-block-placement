package objects

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Tile struct {
	sprite2D  *ebiten.Image
	spriteIso *ebiten.Image
	height    int
}

type TileStack struct {
	stack         []*Tile
	currentHeight int
	maxHeight     int
}

type HalfCube struct {
	*Tile
	height int
}

type Board struct {
	data      [][]*TileStack
	originIso *image.Point
	origin2D  *image.Point
}

var ORIGIN_ISO = image.Point{X: config.ScreenWidth/2 - TILE_WIDTH_ISO/2, Y: config.ScreenHeight/2 - TILE_HEIGHT_ISO/2}
var ORIGIN_2D = image.Point{X: 100, Y: 100}

const (
	TILE_WIDTH_ISO  = 32
	TILE_HEIGHT_ISO = 16
	TILE_WIDTH_2D   = 24
	TILE_HEIGHT_2D  = 18
)

func (ts *TileStack) render2D(screen *ebiten.Image, x int, y int, offset *image.Point) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawX := offset.X + x*TILE_WIDTH_2D
			drawY := offset.Y + y*TILE_HEIGHT_2D
			drawOpts.GeoM.Translate(float64(drawX), float64(drawY))
			screen.DrawImage(tile.sprite2D, drawOpts)
		}
	}
}

func (b *Board) Render2D(screen *ebiten.Image) {
	for y, row := range b.data {
		for x, tileStack := range row {
			tileStack.render2D(screen, x, y, b.origin2D)
		}
	}
}

func (ts *TileStack) renderIso(screen *ebiten.Image, x int, y int, offset *image.Point) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawX := offset.X + (x * (TILE_WIDTH_ISO / 2)) + (y * (TILE_WIDTH_ISO / 2))
			drawY := offset.Y + (y * (TILE_HEIGHT_ISO / 2)) - (x * (TILE_HEIGHT_ISO / 2))
			drawOpts.GeoM.Translate(float64(drawX), float64(drawY))
			screen.DrawImage(tile.spriteIso, drawOpts)
		}
	}
}

func (b *Board) RenderIso(screen *ebiten.Image) {
	for y, row := range b.data {
		for x, tileStack := range row {
			tileStack.renderIso(screen, x, y, b.originIso)
		}
	}
}

func newGroundTile(loader *resource.Loader) *Tile {
	return &Tile{
		sprite2D:  loader.LoadImage(assets.ImgGround2D).Data,
		spriteIso: loader.LoadImage(assets.ImgGroundIso).Data,
		height:    0,
	}
}

func newTileStack(maxHeight int, loader *resource.Loader) *TileStack {
	stack := make([]*Tile, maxHeight)
	stack[0] = newGroundTile(loader)

	return &TileStack{
		stack:         stack,
		currentHeight: stack[0].height,
		maxHeight:     maxHeight,
	}
}

func NewBoard(w int, h int, d int, loader *resource.Loader) *Board {
	data := make([][]*TileStack, w)
	for i := range data {
		data[i] = make([]*TileStack, h)
		for j := range data[i] {
			data[i][j] = newTileStack(d, loader)
		}
	}

	originIso := &image.Point{X: config.ScreenWidth/2 - (w*TILE_WIDTH_ISO)/2, Y: config.ScreenHeight/1.25 - (h*TILE_HEIGHT_ISO)/2}
	origin2D := &image.Point{X: config.ScreenWidth/2 - (w*TILE_WIDTH_2D)/2, Y: config.ScreenHeight/1.5 - (h*TILE_HEIGHT_2D)/2}

	return &Board{
		data:      data,
		originIso: originIso,
		origin2D:  origin2D,
	}
}
