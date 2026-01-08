package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

func createTerminalPanel(state *sys.AppState) fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle("Ghost Terminal", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// TextGrid Log (Placeholder)
	grid := widget.NewTextGrid()
	grid.SetText("Nexus Ops v1.0.0 initialized...\n> System Check: OK\n> Network: CONNECTED\n> Security: ENCRYPTED\n> Waiting for command input...\n_")

	// Create a background for the terminal specifically (darker)
	termBg := canvas.NewRectangle(color.RGBA{R: 10, G: 10, B: 15, A: 250})
	termContainer := container.NewStack(termBg, container.NewPadded(grid))

	// --- Glass Card Styling ---
	content := container.NewBorder(header, nil, nil, nil, termContainer)

	// Background: Semi-transparent dark overlay + White 1px border
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 40, A: 200})
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 50}
	border.StrokeWidth = 1

	card := container.NewStack(bg, border, container.NewPadded(content))

	return container.NewPadded(card)
}
