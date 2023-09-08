package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Point struct {
	X float64
	Y float64
}

type Tile struct {
	sprite2D  *ebiten.Image
	spriteIso *ebiten.Image
	height    int
	pointIso  *Point
	point2D   *Point
}

type TileStack struct {
	stack         []*Tile
	currentHeight int
	maxHeight     int
	surfaceIso    *resolv.Object
	surface2D     *resolv.Object
}

type HalfCube struct {
	*Tile
	height int
}

type Board struct {
	data                [][]*TileStack
	originIso           *Point
	origin2D            *Point
	spaceIso            *resolv.Space
	space2D             *resolv.Space
	cursor              *resolv.Object
	surfaceHighlight2D  *ebiten.Image
	surfaceHighlightIso *ebiten.Image
}

const (
	TILE_WIDTH_ISO  = 32
	TILE_HEIGHT_ISO = 16
	TILE_WIDTH_2D   = 24
	TILE_HEIGHT_2D  = 18
)

func (ts *TileStack) render2D(screen *ebiten.Image) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawOpts.GeoM.Translate(tile.point2D.X, tile.point2D.Y)
			screen.DrawImage(tile.sprite2D, drawOpts)
		}
	}
}

func (b *Board) Render2D(screen *ebiten.Image) {
	for _, row := range b.data {
		for _, tileStack := range row {
			tileStack.render2D(screen)
		}
	}
	x, y := ebiten.CursorPosition()
	hoveredObject := b.space2D.CheckCells(x, y, config.ScreenWidth, config.ScreenHeight)
	if hoveredObject != nil {
		drawOpts := ebiten.DrawImageOptions{}
		drawOpts.GeoM.Translate(hoveredObject.X, hoveredObject.Y)
		screen.DrawImage(b.surfaceHighlight2D, &drawOpts)
	}
}

func (ts *TileStack) renderIso(screen *ebiten.Image) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawOpts.GeoM.Translate(tile.pointIso.X, tile.pointIso.Y)
			screen.DrawImage(tile.spriteIso, drawOpts)
		}
	}
}

func (b *Board) RenderIso(screen *ebiten.Image) {
	for _, row := range b.data {
		for _, tileStack := range row {
			tileStack.renderIso(screen)
		}
	}
	x, y := ebiten.CursorPosition()
	hoveredObject := b.spaceIso.CheckCells(x, y, config.ScreenWidth, config.ScreenHeight)
	if hoveredObject != nil {
		drawOpts := ebiten.DrawImageOptions{}
		drawOpts.GeoM.Translate(hoveredObject.X, hoveredObject.Y)
		screen.DrawImage(b.surfaceHighlightIso, &drawOpts)
	}
}

func newGroundTile(loader *resource.Loader) *Tile {
	return &Tile{
		sprite2D:  loader.LoadImage(assets.ImgGround2D).Data,
		spriteIso: loader.LoadImage(assets.ImgGroundIso).Data,
		height:    0,
	}
}

func newTileStack(x int, y int, maxHeight int, loader *resource.Loader) *TileStack {
	stack := make([]*Tile, maxHeight)
	stack[0] = newGroundTile(loader)

	return &TileStack{
		stack:         stack,
		currentHeight: stack[0].height,
		maxHeight:     maxHeight,
	}
}

func new2DCollision(x float64, y float64) *resolv.Object {
	return resolv.NewObject(x, y, TILE_WIDTH_2D, TILE_HEIGHT_2D)
}

func newIsoCollision(x float64, y float64) *resolv.Object {
	object := resolv.NewObject(x, y, TILE_WIDTH_2D, TILE_HEIGHT_2D)
	object.SetShape(resolv.NewConvexPolygon(
		x, y,
		TILE_WIDTH_ISO/2, 0,
		TILE_WIDTH_ISO, TILE_HEIGHT_ISO/2,
		TILE_WIDTH_ISO/2, TILE_HEIGHT_ISO,
		0, TILE_HEIGHT_ISO/2,
	))
	return object
}

func NewBoard(w int, h int, d int, loader *resource.Loader) *Board {
	data := make([][]*TileStack, w)

	originIso := &Point{
		X: float64(config.ScreenWidth)/2 - float64(w*TILE_WIDTH_ISO)/2,
		Y: float64(config.ScreenHeight)/1.25 - float64(h*TILE_HEIGHT_ISO)/2,
	}
	origin2D := &Point{
		X: float64(config.ScreenWidth)/2 - float64(w*TILE_WIDTH_2D)/2,
		Y: float64(config.ScreenHeight)/1.5 - float64(h*TILE_HEIGHT_2D)/2,
	}
	spaceIso := resolv.NewSpace(config.ScreenWidth, config.ScreenHeight, 1, 1)
	space2D := resolv.NewSpace(config.ScreenWidth, config.ScreenHeight, 1, 1)

	for y := range data {
		data[y] = make([]*TileStack, h)
		for x := range data[y] {
			tileStack := newTileStack(x, y, d, loader)
			tileStack.stack[0].pointIso = &Point{
				X: originIso.X + float64((x*(TILE_WIDTH_ISO/2))+(y*(TILE_WIDTH_ISO/2))),
				Y: originIso.Y + float64((y*(TILE_HEIGHT_ISO/2))-(x*(TILE_HEIGHT_ISO/2))),
			}
			tileStack.surfaceIso = newIsoCollision(tileStack.stack[0].pointIso.X, tileStack.stack[0].pointIso.Y)
			spaceIso.Add(tileStack.surfaceIso)

			tileStack.stack[0].point2D = &Point{
				X: origin2D.X + float64(x*TILE_WIDTH_2D),
				Y: origin2D.Y + float64(y*TILE_HEIGHT_2D),
			}
			tileStack.surface2D = new2DCollision(tileStack.stack[0].point2D.X, tileStack.stack[0].point2D.Y)
			space2D.Add(tileStack.surface2D)

			data[y][x] = tileStack
		}
	}

	surfaceHighlightIso := loader.LoadImage(assets.ImgHoverIso).Data
	surfaceHighlight2D := loader.LoadImage(assets.ImgHover2D).Data

	return &Board{
		data:                data,
		originIso:           originIso,
		origin2D:            origin2D,
		spaceIso:            spaceIso,
		space2D:             space2D,
		surfaceHighlightIso: surfaceHighlightIso,
		surfaceHighlight2D:  surfaceHighlight2D,
	}
}
