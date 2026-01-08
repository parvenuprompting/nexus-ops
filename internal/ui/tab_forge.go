package ui

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/forge"
	"nexus-ops/internal/sys"
)

func createForgePanel(state *sys.AppState, window fyne.Window) fyne.CanvasObject {
	processor := forge.NewProcessor(state)

	// --- Header ---
	label := widget.NewLabelWithStyle("Forge: Image Resizer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	dirLabel := widget.NewLabel("No directory selected")
	var selectedDir string

	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressLabel := widget.NewLabel("Ready")
	progressLabel.Alignment = fyne.TextAlignCenter

	startBtn := widget.NewButton("Start Processing", nil)
	startBtn.Disable()

	selectBtn := widget.NewButton("Select Input Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			selectedDir = uri.Path()
			dirLabel.SetText(fmt.Sprintf("Selected: %s", selectedDir))
			startBtn.Enable()
		}, window)
	})

	startBtn.OnTapped = func() {
		if selectedDir == "" {
			return
		}

		startBtn.Disable()
		selectBtn.Disable()
		progressBar.SetValue(0)
		progressLabel.SetText("Scanning and configuring...")

		// Reset stats
		state.Metrics.TasksCompleted.Store(0)
		state.Metrics.Errors.Store(0)

		// Start in background
		go func() {
			defer func() {
				startBtn.Enable()
				selectBtn.Enable()
			}()

			// How many workers? Let's use runtime.NumCPU()
			workerCount := runtime.NumCPU()
			progressLabel.SetText(fmt.Sprintf("Starting %d workers...", workerCount))

			total, doneChan, err := processor.Start(state.Ctx, selectedDir, workerCount)
			if err != nil {
				dialog.ShowError(err, window)
				progressLabel.SetText(fmt.Sprintf("Error: %v", err))
				return
			}

			if total == 0 {
				progressLabel.SetText("No images found in directory.")
				progressBar.SetValue(1) // Full
				return
			}

			// Watch progress
			progressLabel.SetText(fmt.Sprintf("Processing %d images with %d workers...", total, workerCount))

			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-doneChan:
					// All done
					progressBar.SetValue(1.0)
					progressLabel.SetText(fmt.Sprintf("Processing Complete! (%d images processed)", total))
					return

				case <-state.Ctx.Done():
					progressLabel.SetText("Cancelled.")
					return

				case <-ticker.C:
					done := state.Metrics.TasksCompleted.Load()
					errs := state.Metrics.Errors.Load()
					current := float64(done + errs)

					// Avoid divide by zero if total somehow 0
					if total > 0 {
						progressBar.SetValue(current / float64(total))
					}
				}
			}
		}()
	}

	// --- Glass Card Styling ---
	content := container.NewVBox(
		label,
		layout.NewSpacer(),
		selectBtn,
		dirLabel,
		layout.NewSpacer(),
		startBtn,
		layout.NewSpacer(),
		progressLabel,
		progressBar,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Check Radar for details.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	)

	// Background: Semi-transparent dark overlay + White 1px border
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 40, A: 200})
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 50}
	border.StrokeWidth = 1

	// Card Stack
	card := container.NewStack(bg, border, container.NewPadded(content))

	// Center the card in the tab
	return container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(500, 400)), card))
}
