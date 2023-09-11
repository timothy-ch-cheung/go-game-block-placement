package ui

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
	ViewToggleChangedHandler *widget.CheckboxChangedHandlerFunc
	BlockSizeChangedHandler  *widget.CheckboxChangedHandlerFunc
}

type State struct {
	Renderer       Renderer
	BlockSize      BlockSize
	BlockOperation *BlockOperation
}

type UI struct {
	ebitenUI *ebitenui.UI
	State    *State
}

func (ui *UI) Update() {
	ui.ebitenUI.Update()
}

func (ui *UI) Draw(screen *ebiten.Image) {
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
		loader.LoadImage(assets.ImgSizeBtnFull).Data,
		loader.LoadImage(assets.ImgSizeBtnHalf).Data,
		loader.LoadImage(assets.ImgSizeBtnDisabled).Data,
		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
		}),
	)
}

func newBlockColourRadioBtns(loader *resource.Loader) (*widget.Container, *BlockOperation) {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)
	var checkboxes []*widget.Checkbox
	blockOperation := SELECT

	var cursorBlockChanged widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if int(args.State) > 0 {
			blockOperation = SELECT
		}
	}
	cursorBlock := newCheckbox(
		&cursorBlockChanged,
		loader.LoadImage(assets.ImgCursorBtnSelected).Data,
		loader.LoadImage(assets.ImgCursorBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(cursorBlock)
	checkboxes = append(checkboxes, cursorBlock)

	var blueBlockChanged widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if int(args.State) > 0 {
			blockOperation = PLACE_BLUE
		}
	}
	blueBlock := newCheckbox(
		&blueBlockChanged,
		loader.LoadImage(assets.ImgBlueBlockBtnSelected).Data,
		loader.LoadImage(assets.ImgBlueBlockBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(blueBlock)
	checkboxes = append(checkboxes, blueBlock)

	var redBlockChanged widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if int(args.State) > 0 {
			blockOperation = PLACE_RED
		}
	}
	redBlock := newCheckbox(
		&redBlockChanged,
		loader.LoadImage(assets.ImgRedBlockBtnSelected).Data,
		loader.LoadImage(assets.ImgRedBlockBtnIdle).Data,
		loader.LoadImage(assets.ImgPanelBtnDisabled).Data,
	)
	container.AddChild(redBlock)
	checkboxes = append(checkboxes, redBlock)

	var yellowBlockChanged widget.CheckboxChangedHandlerFunc = func(args *widget.CheckboxChangedEventArgs) {
		if int(args.State) > 0 {
			blockOperation = PLACE_YELLOW
		}
	}
	yellowBlock := newCheckbox(
		&yellowBlockChanged,
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

	radioGroup := widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(elements...),
	)
	radioGroup.SetActive(elements[0])

	return container, &blockOperation
}

func NewUserInterface(handlers *Handlers, loader *resource.Loader) *UI {
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
	viewToggle := newViewToggle(handlers.ViewToggleChangedHandler, loader)
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

	blockSize := HALF
	blockSizeToggle := newSizeToggle(handlers.BlockSizeChangedHandler, loader)
	blockSizeToggle.SetState(widget.WidgetState(blockSize))
	panelContainer.AddChild(blockSizeToggle)

	blockOperationContainer, blockOperation := newBlockColourRadioBtns(loader)
	panelContainer.AddChild(blockOperationContainer)

	rootContainer.AddChild(panelLayout)

	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	state := &State{
		Renderer:       renderer,
		BlockSize:      blockSize,
		BlockOperation: blockOperation,
	}

	return &UI{
		ebitenUI: ui,
		State:    state,
	}
}
