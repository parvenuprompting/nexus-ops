package ui

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

// Setup constructs the UI and attaches it to the window.
func Setup(w fyne.Window, state *sys.AppState) {
	// --- Footer ---
	statusData := binding.NewString()
	statusLabel := widget.NewLabelWithData(statusData)

	// --- Tabs ---
	// Tab 1: Forge
	forgeTab := createForgePanel(state, w)
	// Tab 2: Radar
	radarTab := createRadarPanel(state)

	tabs := container.NewAppTabs(
		container.NewTabItem("Forge", forgeTab),
		container.NewTabItem("Radar", radarTab),
	)

	// --- Main Layout ---
	content := container.NewBorder(nil, statusLabel, nil, nil, tabs)
	w.SetContent(content)

	// --- Global Update Loop ---
	// Updates footer and triggers generic UI refreshes if needed
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

func createRadarPanel(state *sys.AppState) fyne.CanvasObject {
	// We need a list that refreshes.
	// Since Fyne lists are data-driven, we'll keep a local slice of active tasks.

	var dataMu sync.Mutex
	var activeTasks []*sys.Task

	list := widget.NewList(
		func() int {
			dataMu.Lock()
			defer dataMu.Unlock()
			return len(activeTasks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			dataMu.Lock()
			task := activeTasks[i]
			dataMu.Unlock()

			dur := time.Since(task.StartedAt).Round(time.Millisecond)
			o.(*widget.Label).SetText(fmt.Sprintf("[%s] %s (%v)", task.Type, task.ID[:8], dur))
		},
	)

	// Refresher loop for Radar
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-state.Ctx.Done():
				return
			case <-ticker.C:
				// Snapshot active tasks
				var snapshot []*sys.Task
				state.ActiveTasks.Range(func(key, value interface{}) bool {
					if t, ok := value.(*sys.Task); ok {
						snapshot = append(snapshot, t)
					}
					return true
				})

				// Sort by time
				sort.Slice(snapshot, func(i, j int) bool {
					return snapshot[i].StartedAt.Before(snapshot[j].StartedAt)
				})

				dataMu.Lock()
				activeTasks = snapshot
				dataMu.Unlock()

				list.Refresh()
			}
		}
	}()

	return container.NewBorder(nil, nil, nil, nil, list)
}
