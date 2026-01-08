package main

import (
	"context"
	"fmt"
	"os/signal"
	"runtime"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
	"nexus-ops/internal/ui"
)

func main() {
	// 1. AppState & Lifecycle
	// Create a context that listens for SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	metrics := sys.NewMetrics()
	state := &sys.AppState{
		Ctx:     ctx,
		Cancel:  stop,
		Metrics: metrics,
	}

	// 2. Fyne Initialization
	a := app.New()
	a.Settings().SetTheme(ui.NewCyberpunkTheme())
	w := a.NewWindow("Nexus Ops")
	w.Resize(fyne.NewSize(800, 600)) // Increased size for Sidebar layout

	// --- Menu Construction ---
	// File Menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Settings", func() {
			// Placeholder for Settings Overlay
			d := dialog.NewInformation("Settings", "System Configuration Overlay [LOCKED]", w)
			d.Show()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exit", func() {
			a.Quit()
		}),
	)

	// Tools Menu
	toolsMenu := fyne.NewMenu("Tools",
		fyne.NewMenuItem("Force Garbage Collection", func() {
			runtime.GC()
			fmt.Println("Forced Garbage Collection")
		}),
		fyne.NewMenuItem("Reset Metrics", func() {
			state.Metrics.TasksCompleted.Store(0)
			state.Metrics.Errors.Store(0)
			// ActiveGoroutines we shouldn't reset manually as it tracks live state
			fmt.Println("Metrics Reset")
		}),
	)

	// Help Menu
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About Nexus Ops", func() {
			icon := widget.NewIcon(ui.NewCyberpunkTheme().Icon(theme.IconNameInfo))
			content := container.NewHBox(icon, widget.NewLabel("Nexus Hub v1.0.2\nAdvanced Operations Dashboard"))
			dialog.ShowCustom("About", "Close", content, w)
		}),
	)

	mainMenu := fyne.NewMainMenu(
		fileMenu,
		toolsMenu,
		helpMenu,
	)
	w.SetMainMenu(mainMenu)

	// Intercept window close to trigger graceful shutdown
	w.SetOnClosed(func() {
		fmt.Println("Graceful shutdown initiated via window close...")
		state.Cancel() // Cancel the global context
	})

	// 3. UI Setup
	ui.Setup(w, state)

	// 4. Run Loop
	// Handle OS signals that might happen before window closes
	go func() {
		<-ctx.Done()
		fmt.Println("Graceful shutdown initiated via signal...")
		a.Quit() // Quit the Fyne app
	}()

	w.ShowAndRun()
}
