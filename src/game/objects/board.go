package objects

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"
	"github.com/timothy-ch-cheung/go-game-block-placement/ui"

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
	currentIndex  int
	currentHeight int
	maxHeight     int
	isHovered     bool
}

type Board struct {
	data              [][]*TileStack
	objectToTileStack map[string]*TileStack
	originIso         *Point
	origin2D          *Point
	space             *resolv.Space
	cursor            *resolv.Object
}

const (
	TILE_WIDTH_ISO  = 32
	TILE_HEIGHT_ISO = 16
	TILE_WIDTH_2D   = 24
	TILE_HEIGHT_2D  = 18
)

func coordTag(x int, y int) string {
	return fmt.Sprintf("%d%d", x, y)
}

func stackKey(tags []string) string {
	return strings.Join(tags[:], ",")
}

func (ts *TileStack) render2D(screen *ebiten.Image) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawOpts.GeoM.Translate(tile.point2D.X, tile.point2D.Y)
			if ts.isHovered {
				drawOpts.ColorM.RotateHue(1.25)
			}
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
}

func (ts *TileStack) renderIso(screen *ebiten.Image) {
	for _, tile := range ts.stack {
		if tile != nil {
			drawOpts := &ebiten.DrawImageOptions{}
			drawOpts.GeoM.Translate(tile.pointIso.X, tile.pointIso.Y)
			if ts.isHovered {
				drawOpts.ColorM.RotateHue(1.25)
			}
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
}

func (b *Board) Update(renderingMode ui.Renderer) {
	x, y := ebiten.CursorPosition()
	b.cursor.X = float64(x)
	b.cursor.Y = float64(y)
	for _, row := range b.data {
		for _, tileStack := range row {
			tileStack.isHovered = false
		}
	}
	if check := b.cursor.Check(0, 0, "ISO"); check != nil && renderingMode == ui.ISOMETRIC {
		if tileStack := b.objectToTileStack[stackKey(check.Objects[0].Tags())]; tileStack != nil {
			tileStack.isHovered = true
		}
	}
	if check := b.cursor.Check(0, 0, "2D"); check != nil && renderingMode == ui.TWO_DIMENSIONAL {
		if tileStack := b.objectToTileStack[stackKey(check.Objects[0].Tags())]; tileStack != nil {
			tileStack.isHovered = true
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

func newTileStack(x int, y int, maxHeight int, loader *resource.Loader) *TileStack {
	stack := make([]*Tile, maxHeight)
	stack[0] = newGroundTile(loader)

	return &TileStack{
		stack:         stack,
		currentIndex:  0,
		currentHeight: stack[0].height,
		maxHeight:     maxHeight,
	}
}

func new2DCollision(x float64, y float64, tag string) *resolv.Object {
	return resolv.NewObject(x, y, TILE_WIDTH_2D, TILE_HEIGHT_2D, "2D", tag)
}

func newIsoCollision(x float64, y float64, tag string) *resolv.Object {
	object := resolv.NewObject(x, y, TILE_WIDTH_2D, TILE_HEIGHT_2D, "ISO", tag)
	object.SetShape(resolv.NewConvexPolygon(
		x, y,
		TILE_WIDTH_ISO/2, 0,
		TILE_WIDTH_ISO, TILE_HEIGHT_ISO/2,
		TILE_WIDTH_ISO/2, TILE_HEIGHT_ISO,
		0, TILE_HEIGHT_ISO/2,
	))
	return object
}

func NewBoard(w int, h int, d int, cursor *resolv.Object, loader *resource.Loader) *Board {
	data := make([][]*TileStack, w)
	objectToTileStack := make(map[string]*TileStack)

	originIso := &Point{
		X: float64(config.ScreenWidth)/2 - float64(w*TILE_WIDTH_ISO)/2,
		Y: float64(config.ScreenHeight)/1.25 - float64(h*TILE_HEIGHT_ISO)/2,
	}
	origin2D := &Point{
		X: float64(config.ScreenWidth)/2 - float64(w*TILE_WIDTH_2D)/2,
		Y: float64(config.ScreenHeight)/1.5 - float64(h*TILE_HEIGHT_2D)/2,
	}
	space := resolv.NewSpace(config.ScreenWidth, config.ScreenHeight, 1, 1)

	for y := range data {
		data[y] = make([]*TileStack, h)
		for x := range data[y] {
			tileStack := newTileStack(x, y, d, loader)

			tileStack.stack[0].pointIso = &Point{
				X: originIso.X + float64((x*(TILE_WIDTH_ISO/2))+(y*(TILE_WIDTH_ISO/2))),
				Y: originIso.Y + float64((y*(TILE_HEIGHT_ISO/2))-(x*(TILE_HEIGHT_ISO/2))),
			}
			collisionIso := newIsoCollision(tileStack.stack[0].pointIso.X, tileStack.stack[0].pointIso.Y, coordTag(x, y))
			space.Add(collisionIso)
			objectToTileStack[stackKey(collisionIso.Tags())] = tileStack

			tileStack.stack[0].point2D = &Point{
				X: origin2D.X + float64(x*TILE_WIDTH_2D),
				Y: origin2D.Y + float64(y*TILE_HEIGHT_2D),
			}
			collision2D := new2DCollision(tileStack.stack[0].point2D.X, tileStack.stack[0].point2D.Y, coordTag(x, y))
			space.Add(collision2D)
			objectToTileStack[stackKey(collision2D.Tags())] = tileStack

			data[y][x] = tileStack
		}
	}

	space.Add(cursor)

	return &Board{
		data:              data,
		objectToTileStack: objectToTileStack,
		originIso:         originIso,
		origin2D:          origin2D,
		space:             space,
		cursor:            cursor,
	}
}
