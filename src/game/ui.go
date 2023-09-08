package game

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
)

type Handlers struct {
	viewToggleChangedHandler *widget.CheckboxChangedHandlerFunc
}

func newImageNineSlice(img *ebiten.Image, centerWidth int, centerHeight int) *image.NineSlice {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	return image.NewNineSlice(img,
		[3]int{(width - centerWidth) / 2, centerWidth, width - (width-centerWidth)/2 - centerWidth},
		[3]int{(height - centerHeight) / 2, centerHeight, height - (height-centerHeight)/2 - centerHeight})
}

func newViewToggle(handler *widget.CheckboxChangedHandlerFunc, loader *resource.Loader) *widget.Checkbox {
	unchecked := &widget.ButtonImageImage{
		Idle:     loader.LoadImage(assets.ImgViewBtnIso).Data,
		Disabled: loader.LoadImage(assets.ImgViewBtnDisabled).Data,
	}
	checked := &widget.ButtonImageImage{
		Idle:     loader.LoadImage(assets.ImgViewBtn2D).Data,
		Disabled: loader.LoadImage(assets.ImgViewBtnDisabled).Data,
	}
	greyed := &widget.ButtonImageImage{
		Idle:     loader.LoadImage(assets.ImgViewBtnDisabled).Data,
		Disabled: loader.LoadImage(assets.ImgViewBtnDisabled).Data,
	}
	graphic := &widget.CheckboxGraphicImage{
		Unchecked: unchecked,
		Checked:   checked,
		Greyed:    greyed,
	}

	checkboxImg := newImageNineSlice(loader.LoadImage(assets.ImgViewBtnDisabled).Data, 42, 14)
	image := &widget.ButtonImage{
		Idle:         checkboxImg,
		Hover:        checkboxImg,
		Pressed:      checkboxImg,
		PressedHover: checkboxImg,
		Disabled:     checkboxImg,
	}
	return widget.NewCheckbox(
		widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(image)),
		widget.CheckboxOpts.Image(graphic),
		widget.CheckboxOpts.StateChangedHandler(*handler),
	)
}

func newUserInterface(handlers *Handlers, loader *resource.Loader) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}))))

	rootContainer.AddChild(newViewToggle(handlers.viewToggleChangedHandler, loader))

	return &ebitenui.UI{
		Container: rootContainer,
	}
}
