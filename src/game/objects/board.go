package objects

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
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
	height    ui.BlockSize
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
	loader            *resource.Loader
	maxHeight         int
}

const (
	TILE_WIDTH_ISO      = 32
	TILE_HEIGHT_ISO     = 16
	TILE_FULL_DEPTH_ISO = 17
	TILE_HALF_DEPTH_ISO = 9
	TILE_WIDTH_2D       = 24
	TILE_HEIGHT_2D      = 18
	TILE_FULL_DEPTH_2D  = 14
	TILE_HALF_DEPTH_2D  = 7
)

func coordTag(x int, y int) string {
	return fmt.Sprintf("%d-%d", x, y)
}

func stackKey(tags []string) string {
	return strings.Join(tags[:], ",")
}

func (ts *TileStack) addTile(blockSize ui.BlockSize, blockOperation ui.BlockOperation, loader *resource.Loader) {
	currentBlock := ts.stack[ts.currentIndex]
	newBlock := newBlockTile(blockSize, blockOperation, loader)

	yIncrementIso := TILE_FULL_DEPTH_ISO
	yIncrement2D := TILE_FULL_DEPTH_2D
	if blockSize == ui.HALF {
		yIncrementIso = TILE_HALF_DEPTH_ISO
		yIncrement2D = TILE_HALF_DEPTH_2D
	}

	newBlock.pointIso = &Point{X: currentBlock.pointIso.X, Y: currentBlock.pointIso.Y - float64(yIncrementIso)}
	newBlock.point2D = &Point{X: currentBlock.point2D.X, Y: currentBlock.point2D.Y - float64(yIncrement2D)}
	ts.stack = append(ts.stack, newBlock)
	ts.currentHeight += newBlock.height.GetHeight()
	ts.currentIndex += 1
}

func (ts *TileStack) deleteTopTile() {
	if ts.currentIndex < 1 {
		return
	}
	ts.currentHeight -= ts.stack[ts.currentIndex].height.GetHeight()
	ts.stack = ts.stack[:len(ts.stack)-1]
	ts.currentIndex -= 1

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
	for j := 0; j < len(b.data); j++ {
		for i := len(b.data[j]) - 1; i >= 0; i-- {
			b.data[j][i].renderIso(screen)
		}
	}
}

func (b *Board) canPlaceBlock(tileStack *TileStack, state *ui.State) bool {
	size := state.BlockSize.GetHeight()
	return tileStack.currentHeight+size <= b.maxHeight && *state.BlockOperation != ui.SELECT
}

func (b *Board) Update(state *ui.State, handler *input.Handler) {
	x, y := ebiten.CursorPosition()
	b.cursor.X = float64(x)
	b.cursor.Y = float64(y)
	for _, row := range b.data {
		for _, tileStack := range row {
			tileStack.isHovered = false
		}
	}

	if check := b.cursor.Check(0, 0, "ISO"); check != nil && state.Renderer == ui.ISOMETRIC {
		if tileStack := b.objectToTileStack[stackKey(check.Objects[0].Tags())]; tileStack != nil {
			tileStack.isHovered = true
			if handler.ActionIsJustPressed(ui.ActionSelect) && b.canPlaceBlock(tileStack, state) {
				tileStack.addTile(state.BlockSize, *state.BlockOperation, b.loader)
			} else if handler.ActionIsJustPressed(ui.ActionDelete) {
				tileStack.deleteTopTile()
			}
		}
	}
	if check := b.cursor.Check(0, 0, "2D"); check != nil && state.Renderer == ui.TWO_DIMENSIONAL {
		if tileStack := b.objectToTileStack[stackKey(check.Objects[0].Tags())]; tileStack != nil {
			tileStack.isHovered = true
			if handler.ActionIsJustPressed(ui.ActionSelect) && b.canPlaceBlock(tileStack, state) {
				tileStack.addTile(state.BlockSize, *state.BlockOperation, b.loader)
			} else if handler.ActionIsJustPressed(ui.ActionDelete) {
				tileStack.deleteTopTile()
			}
		}
	}
}

func newGroundTile(loader *resource.Loader) *Tile {
	return &Tile{
		sprite2D:  loader.LoadImage(assets.ImgGround2D).Data,
		spriteIso: loader.LoadImage(assets.ImgGroundIso).Data,
		height:    ui.FLAT,
	}
}

func newBlockTile(blockSize ui.BlockSize, blockOperation ui.BlockOperation, loader *resource.Loader) *Tile {
	var sprite2D *ebiten.Image
	var spriteIso *ebiten.Image
	switch true {
	case (blockOperation == ui.PLACE_BLUE && blockSize == ui.HALF):
		sprite2D = loader.LoadImage(assets.ImgBlueHalfCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgBlueHalfCubeIso).Data

	case (blockOperation == ui.PLACE_RED && blockSize == ui.HALF):
		sprite2D = loader.LoadImage(assets.ImgRedHalfCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgRedHalfCubeIso).Data

	case (blockOperation == ui.PLACE_YELLOW && blockSize == ui.HALF):
		sprite2D = loader.LoadImage(assets.ImgYellowHalfCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgYellowHalfCubeIso).Data

	case (blockOperation == ui.PLACE_BLUE && blockSize == ui.FULL):
		sprite2D = loader.LoadImage(assets.ImgBlueCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgBlueCubeIso).Data

	case (blockOperation == ui.PLACE_RED && blockSize == ui.FULL):
		sprite2D = loader.LoadImage(assets.ImgRedCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgRedCubeIso).Data

	case (blockOperation == ui.PLACE_YELLOW && blockSize == ui.FULL):
		sprite2D = loader.LoadImage(assets.ImgYellowCube2D).Data
		spriteIso = loader.LoadImage(assets.ImgYellowCubeIso).Data
	}

	return &Tile{
		sprite2D:  sprite2D,
		spriteIso: spriteIso,
		height:    blockSize,
	}
}

func newTileStack(x int, y int, maxHeight int, loader *resource.Loader) *TileStack {
	stack := make([]*Tile, 1, maxHeight)
	stack[0] = newGroundTile(loader)

	return &TileStack{
		stack:         stack,
		currentIndex:  0,
		currentHeight: stack[0].height.GetHeight(),
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
		Y: float64(config.ScreenHeight)/1.75 - float64(h*TILE_HEIGHT_2D)/2,
	}
	space := resolv.NewSpace(config.ScreenWidth, config.ScreenHeight, 1, 1)

	for y := range data {
		data[y] = make([]*TileStack, h)
		for x := range data[y] {
			tileStack := newTileStack(x, y, d, loader)

			xIso, yIso := calculateIsoCoord(originIso, x, y)
			tileStack.stack[0].pointIso = &Point{X: xIso, Y: yIso}
			collisionIso := newIsoCollision(tileStack.stack[0].pointIso.X, tileStack.stack[0].pointIso.Y, coordTag(x, y))
			space.Add(collisionIso)
			objectToTileStack[stackKey(collisionIso.Tags())] = tileStack

			x2D, y2D := calculate2DCoord(origin2D, x, y)
			tileStack.stack[0].point2D = &Point{X: x2D, Y: y2D}
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
		loader:            loader,
		maxHeight:         d * 2,
	}
}

func calculateIsoCoord(originIso *Point, x int, y int) (float64, float64) {
	xIso := originIso.X + float64((x*(TILE_WIDTH_ISO/2))+(y*(TILE_WIDTH_ISO/2)))
	yIso := originIso.Y + float64((y*(TILE_HEIGHT_ISO/2))-(x*(TILE_HEIGHT_ISO/2)))
	return xIso, yIso
}

func calculate2DCoord(origin2D *Point, x int, y int) (float64, float64) {
	x2D := origin2D.X + float64(x*TILE_WIDTH_2D)
	y2D := origin2D.Y + float64(y*TILE_HEIGHT_2D)
	return x2D, y2D
}
