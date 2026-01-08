package ui

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/forge"
	"nexus-ops/internal/sys"
)

func createForgePanel(state *sys.AppState, window fyne.Window) fyne.CanvasObject {
	processor := forge.NewProcessor(state)

	label := widget.NewLabel("Forge: Image Resizer")

	dirLabel := widget.NewLabel("No directory selected")
	var selectedDir string

	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressLabel := widget.NewLabel("Ready")

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

			total, err := processor.Start(state.Ctx, selectedDir, workerCount)
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

			// We can poll metrics to update progress bar smoothly
			// Simple loop until done
			// Note: Start() returns immediately after queueing, but the actual work happens async.
			// Start() logic in previous step: Enqueues then returns. The workers are running.
			// We need to know when they are finished.
			// Start() returned total count.
			// We can check TasksCompleted + Errors == Total.

			for {
				done := state.Metrics.TasksCompleted.Load()
				errs := state.Metrics.Errors.Load()
				current := float64(done + errs)

				progressBar.SetValue(current / float64(total))

				if current >= float64(total) {
					break
				}

				if state.Ctx.Err() != nil {
					progressLabel.SetText("Cancelled.")
					return
				}

				// Optional: sleep slightly to not spam lock
				// time.Sleep(100 * time.Millisecond) // Don't block Fyne thread (we are in goroutine though)
				// Using internal ticker or just busy loop with small yield
			}
			progressLabel.SetText("Processing Complete!")
		}()
	}

	return container.NewVBox(
		label,
		selectBtn,
		dirLabel,
		widget.NewSeparator(),
		startBtn,
		progressLabel,
		progressBar,
		widget.NewSeparator(),
		widget.NewLabel("Check Radar for details."),
	)
}
