package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

// Setup constructs the UI and attaches it to the window.
func Setup(w fyne.Window, state *sys.AppState) {
	// --- Footer ---
	statusData := binding.NewString()
	statusLabel := widget.NewLabelWithData(statusData)

	// --- Modules ---
	forgePanel := createForgePanel(state, w)
	radarPanel := createRadarPanel(state)
	vaultPanel := createVaultPanel(state)
	siphonPanel := createSiphonPanel(state)
	termPanel := createTerminalPanel(state)

	// --- Navigation Logic ---
	// We'll use a container stack to swap content
	contentStack := container.NewStack()
	// Default View
	contentStack.Add(forgePanel)

	// We utilize a simple map or switch to swap content.
	// Sidebar Buttons

	// Helper to create nav buttons
	navBtn := func(label string, panel fyne.CanvasObject) *widget.Button {
		return widget.NewButton(label, func() {
			contentStack.Objects = []fyne.CanvasObject{panel}
			contentStack.Refresh()
		})
	}

	sidebar := container.NewVBox(
		widget.NewLabelWithStyle("NEXUS HUB", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		navBtn("Forge", forgePanel),
		navBtn("Radar", radarPanel),
		navBtn("The Vault", vaultPanel),
		navBtn("Siphon", siphonPanel),
		navBtn("Terminal", termPanel),
		layout.NewSpacer(), // Push content up
		widget.NewLabelWithStyle("v1.0.2", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	)

	// --- Layout: Sidebar + Content ---
	// Use HSplit for resizable sidebar
	split := container.NewHSplit(
		container.NewPadded(sidebar),
		container.NewPadded(contentStack),
	)
	split.SetOffset(0.2) // 20% width for sidebar

	// --- Main Layout ---
	// Footer at bottom
	content := container.NewBorder(nil, statusLabel, nil, nil, split)
	w.SetContent(content)

	// --- Global Update Loop ---
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-state.Ctx.Done():
				return
			case <-ticker.C:
				active := state.Metrics.ActiveGoroutines.Load()
				completed := state.Metrics.TasksCompleted.Load()
				errs := state.Metrics.Errors.Load()
				statusData.Set(fmt.Sprintf("Active: %d | Completed: %d | Errors: %d", active, completed, errs))
			}
		}
	}()
}
