package main

import (
	"context"
	"fmt"
	"nexus-ops/internal/services"
	"nexus-ops/internal/sys"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	state         *sys.AppState
	opsService    *services.OpsService
	vaultService  *services.VaultService
	siphonService *services.SiphonService
}

// NewApp creates a new App application struct
func NewApp(state *sys.AppState) *App {
	return &App{
		state:         state,
		opsService:    services.NewOpsService(state),
		vaultService:  services.NewVaultService(),
		siphonService: services.NewSiphonService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start the update listener to push events to frontend
	go a.monitorUpdates()
}

// monitorUpdates listens to the AppState.UpdateChan and emits events
func (a *App) monitorUpdates() {
	for {
		select {
		case <-a.state.Ctx.Done():
			return
		case <-a.state.UpdateChan:
			// Push Active Tasks update
			tasks := a.GetActiveTasks()
			runtime.EventsEmit(a.ctx, "activeTasks:update", tasks)

			// Push Metrics update
			metrics := a.GetMetrics()
			runtime.EventsEmit(a.ctx, "metrics:update", metrics)
		}
	}
}

// -- Exposed Methods --

// SelectDirectory opens a dialog and returns the path
func (a *App) SelectDirectory() string {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Images Folder",
	})
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error selecting directory: %v", err)
		return ""
	}
	return path
}

// StartProcessing initiates the forge processor
// Returns 0 immediately, events will track progress
func (a *App) StartProcessing(dir string, settings sys.ForgeSettings) int {
	err := a.opsService.StartProcessing(a.ctx, dir, settings)
	if err != nil {
		appErr := sys.NewError("ERR_SCAN_FAILED", err.Error())
		runtime.EventsEmit(a.ctx, "processing:error", appErr)
		return 0
	}
	return 0 // Return immediately
}

// GetMetrics returns the current metrics
func (a *App) GetMetrics() sys.MetricsDTO {
	return sys.MetricsDTO{
		ActiveGoroutines: a.state.Metrics.ActiveGoroutines.Load(),
		TasksCompleted:   a.state.Metrics.TasksCompleted.Load(),
		Errors:           a.state.Metrics.Errors.Load(),
	}
}

// GetActiveTasks returns the current active tasks
func (a *App) GetActiveTasks() []*sys.Task {
	return a.state.SnapshotActiveTasks()
}

// RunCommand executes a terminal command
func (a *App) RunCommand(cmd string) string {
	output, err := a.opsService.RunTerminalCommand(a.ctx, cmd)
	if err != nil {
		// Return output + error message if command failed (e.g. stderr)
		return fmt.Sprintf("%s\nError: %s", output, err.Error())
	}
	return output
}

// -- Vault Methods --

func (a *App) GetSecrets() []services.Secret {
	secrets, err := a.vaultService.LoadSecrets()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to load secrets: %v", err)
		return []services.Secret{}
	}
	return secrets
}

func (a *App) AddSecret(key, value string) bool {
	if err := a.vaultService.AddSecret(key, value); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to add secret: %v", err)
		return false
	}
	return true
}

func (a *App) RemoveSecret(key string) bool {
	if err := a.vaultService.RemoveSecret(key); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to remove secret: %v", err)
		return false
	}
	return true
}

// -- Siphon Methods --

type SiphonResult struct {
	StatusCode int    `json:"statusCode"`
	LatencyMs  int64  `json:"latencyMs"`
	Error      string `json:"error"`
}

func (a *App) SiphonCheckHTTP(url string) SiphonResult {
	code, latency, err := a.siphonService.CheckHTTP(url)
	if err != nil {
		return SiphonResult{Error: err.Error()}
	}
	return SiphonResult{StatusCode: code, LatencyMs: latency}
}
