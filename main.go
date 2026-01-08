package main

import (
	"context"
	"embed"
	"os/signal"
	"syscall"

	"nexus-ops/internal/sys"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 1. AppState & Lifecycle
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	metrics := sys.NewMetrics()
	state := &sys.AppState{
		Ctx:     ctx,
		Cancel:  stop,
		Metrics: metrics,
	}

	// 2. Create App instance (Bridge)
	app := NewApp(state)

	// 3. Wails Run
	err := wails.Run(&options.App{
		Title:  "Nexus Ops",
		Width:  1024, // Required Config
		Height: 768,  // Required Config
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 20, B: 30, A: 255}, // Dark background
		Frameless:        true,                                       // Required Config
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			Theme:                windows.SystemDefault,
			// Custom backdrop?
			BackdropType: windows.Mica, // For that nice blur if supported
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
