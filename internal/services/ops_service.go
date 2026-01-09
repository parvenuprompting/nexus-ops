package services

import (
	"context"
	"fmt"
	"os/exec"
	nativeRuntime "runtime"
	"time"

	"nexus-ops/internal/forge"
	"nexus-ops/internal/sys"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// OpsService handles core operation logic
type OpsService struct {
	state *sys.AppState
}

// NewOpsService creates a new OpsService
func NewOpsService(state *sys.AppState) *OpsService {
	return &OpsService{
		state: state,
	}
}

// StartProcessing initiates the forge processor
func (s *OpsService) StartProcessing(ctx context.Context, dir string, settings sys.ForgeSettings) error {
	// Reset metrics
	s.state.Metrics.TasksCompleted.Store(0)
	s.state.Metrics.Errors.Store(0)

	workerCount := nativeRuntime.NumCPU()
	processor := forge.NewProcessor(s.state)

	// Start returns (scanFoundChan, workersDoneChan, err)
	scanFoundChan, workersDoneChan, err := processor.Start(s.state.Ctx, dir, workerCount, settings)
	if err != nil {
		return err
	}

	// Monitor progress and completion in background
	// Monitor progress and completion in background
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log panic and try to signal completion to avoid hanging UI
				runtime.LogErrorf(ctx, "PANIC in OpsService Monitor: %v", r)
				runtime.EventsEmit(ctx, "processing:complete", true)
			}
		}()

		// 1. Wait for Scan Result (Streaming)
		total := 0
		select {
		case t := <-scanFoundChan:
			total = t
			runtime.EventsEmit(ctx, "processing:started", total) // Signals Switch to Processing State
		case <-s.state.Ctx.Done():
			return
		}

		// 2. Monitor Workers
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-workersDoneChan:
				runtime.EventsEmit(ctx, "processing:complete", true)
				return
			case <-s.state.Ctx.Done():
				return
			case <-ticker.C:
				// Optional heartbeat
			}
		}
	}()

	return nil
}

// RunTerminalCommand executes a shell command and returns output
func (s *OpsService) RunTerminalCommand(ctx context.Context, cmdStr string) (string, error) {
	// Security Note: This is a basic implementation.
	// In production, you'd want to validate inputs or use a PTY library.
	var cmd *exec.Cmd

	if nativeRuntime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", cmdStr)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}
