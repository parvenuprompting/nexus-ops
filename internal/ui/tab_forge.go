package ui

import (
	"fmt"
	"image/color"
	"math"
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
		progressLabel.SetText("Initializing Scanner...")

		// Reset stats
		state.Metrics.TasksCompleted.Store(0)
		state.Metrics.Errors.Store(0)

		// Start in background
		go func() {
			defer func() {
				// Re-enable UI on main thread inside callback or just rely on Fyne's thread safety for simple Enables?
				// Fyne widgets are generally thread-safe for basic Set calls, but let's be safe.
				// However, defer runs at end of goroutine.
				startBtn.Enable()
				selectBtn.Enable()
			}()

			workerCount := runtime.NumCPU()

			// Phase 1: Streaming Start
			scanDoneChan, workersDoneChan, err := processor.Start(state.Ctx, selectedDir, workerCount)
			if err != nil {
				dialog.ShowError(err, window)
				progressLabel.SetText(fmt.Sprintf("Error: %v", err))
				return
			}

			// Progress Loop
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			scanning := true
			total := 0
			progressLabel.SetText("Scanning... (0 found)")

			for {
				select {
				case <-state.Ctx.Done():
					progressLabel.SetText("Cancelled.")
					return

				case t, ok := <-scanDoneChan:
					if ok {
						total = t
						scanning = false
						if total == 0 {
							progressLabel.SetText("No images found.")
							progressBar.SetValue(1.0)
							// We still wait for workersDoneChan which should close immediately
						}
					}

				case <-workersDoneChan:
					progressBar.SetValue(1.0)
					progressLabel.SetText(fmt.Sprintf("Complete! %d images processed.", total))
					return

				case <-ticker.C:
					// Update UI
					done := state.Metrics.TasksCompleted.Load()
					errs := state.Metrics.Errors.Load()

					if scanning {
						// Indeterminate Pulse effect
						val := progressBar.Value + 0.05
						if val > 1.0 {
							val = 0
						}
						progressBar.SetValue(val)

						// Try to estimate found based on processed if we don't have total?
						// Actually scanning is fast, but let's show some activity.
						// We don't have "found so far" count from scanner unless we added another channel.
						// "Scanning..." is sufficient.
					} else {
						// Deterministic
						if total > 0 {
							current := float64(done + errs)
							progressBar.SetValue(math.Min(current/float64(total), 1.0))
							progressLabel.SetText(fmt.Sprintf("Processing... (%d/%d)", int(done+errs), total))
						}
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

	// Background
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 40, A: 200})
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 50}
	border.StrokeWidth = 1

	card := container.NewStack(bg, border, container.NewPadded(content))

	return container.NewCenter(container.New(layout.NewGridWrapLayout(fyne.NewSize(500, 400)), card))
}
