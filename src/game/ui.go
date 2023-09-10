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

type Renderer int

const (
	ISOMETRIC Renderer = iota
	TWO_DIMENSIONAL
)

type BlockSize int

const (
	HALF BlockSize = iota
	FULL
)

type Handlers struct {
	viewToggleChangedHandler *widget.CheckboxChangedHandlerFunc
	blockSizeChangedHandler  *widget.CheckboxChangedHandlerFunc
}

type UI struct {
	ebitenUI  *ebitenui.UI
	renderer  Renderer
	blockSize BlockSize
}

func (ui *UI) update() {
	ui.ebitenUI.Update()
}

func (ui *UI) draw(screen *ebiten.Image) {
	ui.ebitenUI.Draw(screen)
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

	checkboxImg := newImageNineSlice(btnDisabledImg, 16, 16)
	image := &widget.ButtonImage{
		Idle:         checkboxImg,
		Hover:        checkboxImg,
		Pressed:      checkboxImg,
		PressedHover: checkboxImg,
		Disabled:     checkboxImg,
	}

	if handler == nil {
		return widget.NewCheckbox(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(image)),
			widget.CheckboxOpts.Image(graphic),
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.WidgetOpts(
				opts...,
			)),
		)
	} else {
		return widget.NewCheckbox(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(image)),
			widget.CheckboxOpts.Image(graphic),
			widget.CheckboxOpts.StateChangedHandler(*handler),
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.WidgetOpts(
				opts...,
			)),
		)
	}
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

func newBlockColourRadioBtns(loader *resource.Loader) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)
	var checkboxes []*widget.Checkbox

	cursorBlock := newCheckbox(nil,
		loader.LoadImage(assets.ImgCursorBtnSelected).Data,
		loader.LoadImage(assets.ImgCursorBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(cursorBlock)
	checkboxes = append(checkboxes, cursorBlock)

	blueBlock := newCheckbox(nil,
		loader.LoadImage(assets.ImgBlueBlockBtnSelected).Data,
		loader.LoadImage(assets.ImgBlueBlockBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(blueBlock)
	checkboxes = append(checkboxes, blueBlock)

	redBlock := newCheckbox(nil,
		loader.LoadImage(assets.ImgRedBlockBtnSelected).Data,
		loader.LoadImage(assets.ImgRedBlockBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(redBlock)
	checkboxes = append(checkboxes, redBlock)

	yellowBlock := newCheckbox(nil,
		loader.LoadImage(assets.ImgYellowBlockBtnSelected).Data,
		loader.LoadImage(assets.ImgYellowBlockBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(yellowBlock)
	checkboxes = append(checkboxes, yellowBlock)

	elements := []widget.RadioGroupElement{}
	for _, cb := range checkboxes {
		elements = append(elements, cb)
	}

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(elements...),
	)
	return container
}

func newUserInterface(handlers *Handlers, loader *resource.Loader) *UI {
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

	renderer := ISOMETRIC
	viewContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
	)
	viewToggle := newViewToggle(handlers.viewToggleChangedHandler, loader)
	viewToggle.SetState(widget.WidgetState(renderer))
	viewContainer.AddChild(viewToggle)
	rootContainer.AddChild(viewContainer)

	panelLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
	)
	panelContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(20),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchVertical: true,
		})),
	)
	panelLayout.AddChild(panelContainer)

	blockSize := FULL
	blockSizeToggle := newSizeToggle(handlers.blockSizeChangedHandler, loader)
	blockSizeToggle.SetState(widget.WidgetState(blockSize))
	panelContainer.AddChild(blockSizeToggle)
	panelContainer.AddChild(newBlockColourRadioBtns(loader))
	rootContainer.AddChild(panelLayout)

	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	return &UI{
		ebitenUI:  ui,
		renderer:  renderer,
		blockSize: blockSize,
	}
}
