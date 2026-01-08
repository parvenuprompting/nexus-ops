package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

func createSiphonPanel(state *sys.AppState) fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle("Siphon: Data Stream", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// Placeholder Content
	info := widget.NewLabel("Waiting for active data streams...\n\n[Status: IDLE]\n[Bandwidth: 0 kbps]")
	info.Alignment = fyne.TextAlignCenter

	// --- Glass Card Styling ---
	content := container.NewVBox(
		header,
		layout.NewSpacer(),
		info,
		layout.NewSpacer(),
	)

	// Background: Semi-transparent dark overlay + White 1px border
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 40, A: 200})
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 50}
	border.StrokeWidth = 1

	card := container.NewStack(bg, border, container.NewPadded(content))

	return container.NewPadded(card)
}
