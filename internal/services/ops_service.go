package services

import (
	"context"
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
func (s *OpsService) StartProcessing(ctx context.Context, dir string) error {
	// Reset metrics
	s.state.Metrics.TasksCompleted.Store(0)
	s.state.Metrics.Errors.Store(0)

	workerCount := nativeRuntime.NumCPU()
	processor := forge.NewProcessor(s.state)

	// Start returns (scanFoundChan, workersDoneChan, err)
	scanFoundChan, workersDoneChan, err := processor.Start(s.state.Ctx, dir, workerCount)
	if err != nil {
		return err
	}

	// Monitor progress and completion in background
	go func() {
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
