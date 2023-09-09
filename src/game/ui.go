package game

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
	"github.com/timothy-ch-cheung/go-game-block-placement/game/config"
)

type Handlers struct {
	viewToggleChangedHandler *widget.CheckboxChangedHandlerFunc
	blockSizeChangedHandler  *widget.CheckboxChangedHandlerFunc
}

func newImageNineSlice(img *ebiten.Image, centerWidth int, centerHeight int) *image.NineSlice {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	return image.NewNineSlice(img,
		[3]int{(width - centerWidth) / 2, centerWidth, width - (width-centerWidth)/2 - centerWidth},
		[3]int{(height - centerHeight) / 2, centerHeight, height - (height-centerHeight)/2 - centerHeight})
}

func newCheckbox(
	handler *widget.CheckboxChangedHandlerFunc,
	btnCheckedImg *ebiten.Image,
	btnUncheckedImg *ebiten.Image,
	btnDisabledImg *ebiten.Image,
	opts ...widget.WidgetOpt,
) *widget.Checkbox {

	unchecked := &widget.ButtonImageImage{
		Idle:     btnUncheckedImg,
		Disabled: btnDisabledImg,
	}
	checked := &widget.ButtonImageImage{
		Idle:     btnCheckedImg,
		Disabled: btnDisabledImg,
	}
	greyed := &widget.ButtonImageImage{
		Idle:     btnDisabledImg,
		Disabled: btnDisabledImg,
	}
	graphic := &widget.CheckboxGraphicImage{
		Unchecked: unchecked,
		Checked:   checked,
		Greyed:    greyed,
	}

	checkboxImg := newImageNineSlice(btnDisabledImg, 42, 14)
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
		widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.WidgetOpts(
			opts...,
		)),
	)
}

func newViewToggle(handler *widget.CheckboxChangedHandlerFunc, loader *resource.Loader) *widget.Checkbox {
	return newCheckbox(
		handler,
		loader.LoadImage(assets.ImgViewBtn2D).Data,
		loader.LoadImage(assets.ImgViewBtnIso).Data,
		loader.LoadImage(assets.ImgViewBtnDisabled).Data,
		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionStart,
			HorizontalPosition: widget.AnchorLayoutPositionStart,
		}),
	)
}

func newSizeToggle(handler *widget.CheckboxChangedHandlerFunc, loader *resource.Loader) *widget.Checkbox {
	return newCheckbox(
		handler,
		loader.LoadImage(assets.ImgSizeBtnHalf).Data,
		loader.LoadImage(assets.ImgSizeBtnFull).Data,
		loader.LoadImage(assets.ImgSizeBtnDisabled).Data,
		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
		}),
	)
}

func newUserInterface(handlers *Handlers, loader *resource.Loader) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(config.ScreenHeight*0.85),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    5,
					Bottom: 5,
					Left:   5,
					Right:  5,
				}),
			),
		),
	)

	viewContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
	)
	viewContainer.AddChild(newViewToggle(handlers.viewToggleChangedHandler, loader))
	rootContainer.AddChild(viewContainer)

	panelContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
	)
	panelContainer.AddChild(newSizeToggle(handlers.blockSizeChangedHandler, loader))
	rootContainer.AddChild(panelContainer)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}
