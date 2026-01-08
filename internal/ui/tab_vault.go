package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

func createVaultPanel(state *sys.AppState) fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle("The Vault: Encrypted Assets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// Data Table (Placeholder)
	data := [][]string{
		{"ID", "Asset Name", "Encryption", "Status"},
		{"0x1A", "Database_Creds.kdbx", "AES-256", "LOCKED"},
		{"0x2B", "API_Keys_Prod.env", "ChaCha20", "LOCKED"},
		{"0x3C", "User_Backups_2025.tar.gz", "AES-256", "SYNCED"},
		{"0x4D", "System_Logs_Audit.log", "None", "OPEN"},
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("Template", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		},
	)

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)

	// --- Glass Card Styling ---
	content := container.NewBorder(header, nil, nil, nil, table)

	// Background: Semi-transparent dark overlay + White 1px border
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 40, A: 200})
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 50}
	border.StrokeWidth = 1

	card := container.NewStack(bg, border, container.NewPadded(content))

	return container.NewPadded(card)
}
