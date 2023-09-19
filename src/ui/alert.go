package ui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/timothy-ch-cheung/go-game-block-placement/assets"
)

const (
	MAX_TICK     = 120
	ANIMATE_TICK = 20
)

type AlertText struct {
	widget      *widget.Text
	color       color.Color
	currentTick int
	isVisible   bool
}

func newAlertText(text string, textColor color.Color, loader *resource.Loader) *AlertText {
	alertWidget := widget.NewText(widget.TextOpts.Text(text, loader.LoadFont(assets.FontDefault).Face, textColor))
	return &AlertText{
		widget: alertWidget,
		color:  textColor,
	}
}

func (alert *AlertText) Animate() {
	alert.isVisible = true
	alert.currentTick = 0
	alert.widget.Color = alert.color
}

func (alert *AlertText) update() {
	if alert.isVisible && alert.currentTick < MAX_TICK {
		if alert.currentTick%ANIMATE_TICK == 0 {
			if alert.widget.Color == alert.color {
				alert.widget.Color = color.Transparent
			} else {
				alert.widget.Color = alert.color
			}
		}
		alert.currentTick++
	} else {
		alert.isVisible = false
		alert.widget.Color = color.Transparent
	}
}
