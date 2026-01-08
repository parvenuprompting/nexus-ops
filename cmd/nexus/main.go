package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

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
	w.Resize(fyne.NewSize(600, 400))

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
