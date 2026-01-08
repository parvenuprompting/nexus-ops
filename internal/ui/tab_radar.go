package ui

import (
	"fmt"
	"image/color"
	"sort"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"nexus-ops/internal/sys"
)

func createRadarPanel(state *sys.AppState) fyne.CanvasObject {
	var dataMu sync.Mutex
	var activeTasks []*sys.Task

	// Empty State Container
	emptyLabel := widget.NewLabelWithStyle("No active tasks. System idle.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	emptyIcon := widget.NewIcon(nil) // Could use a theme icon if available, or just text
	emptyState := container.NewCenter(container.NewVBox(emptyIcon, emptyLabel))
	emptyState.Hide()

	// List
	list := widget.NewList(
		func() int {
			dataMu.Lock()
			defer dataMu.Unlock()
			return len(activeTasks)
		},
		func() fyne.CanvasObject {
			// Template Item: Badge + ID + Duration
			badge := canvas.NewRectangle(color.Transparent)
			badge.SetMinSize(fyne.NewSize(10, 10))

			badgeLabel := widget.NewLabelWithStyle("TYPE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			badgeLabel.TextStyle.Monospace = true

			// Compact container for badge
			badgeContainer := container.NewMax(badge, container.NewCenter(badgeLabel))

			// Main text
			label := widget.NewLabel("Task ID")
			label.TextStyle.Monospace = true

			return container.NewBorder(nil, nil, badgeContainer, nil, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			dataMu.Lock()
			if i >= len(activeTasks) {
				dataMu.Unlock()
				return
			}
			task := activeTasks[i]
			dataMu.Unlock()

			// Update Content
			border := o.(*fyne.Container)
			badgeContainer := border.Objects[1].(*fyne.Container)
			badgeBg := badgeContainer.Objects[0].(*canvas.Rectangle)
			badgeLbl := badgeContainer.Objects[1].(*fyne.Container).Objects[0].(*widget.Label) // Center -> Label

			mainLbl := border.Objects[0].(*widget.Label)

			// Badge Color based on type
			switch task.Type {
			case sys.TaskTypeCPU:
				badgeBg.FillColor = color.RGBA{R: 255, G: 100, B: 0, A: 200} // Orange
				badgeLbl.SetText("CPU")
			case sys.TaskTypeIO:
				badgeBg.FillColor = color.RGBA{R: 0, G: 200, B: 100, A: 200} // Green
				badgeLbl.SetText("I/O")
			case sys.TaskTypeImageResize:
				badgeBg.FillColor = color.RGBA{R: 200, G: 0, B: 255, A: 200} // Purple/Pink
				badgeLbl.SetText("IMG")
			default:
				badgeBg.FillColor = color.RGBA{R: 100, G: 100, B: 100, A: 200} // Grey
				badgeLbl.SetText("???")
			}

			dur := time.Since(task.StartedAt).Round(time.Millisecond)
			mainLbl.SetText(fmt.Sprintf("%s (%v)", task.ID[:12], dur))
		},
	)

	// Wrap list in scroll container/border
	listContainer := container.NewBorder(nil, nil, nil, nil, list)

	// Main container stack (List + Empty State)
	mainStack := container.NewStack(listContainer, emptyState)

	updateList := func() {
		// Snapshot Active Tasks
		snapshot := state.SnapshotActiveTasks()

		sort.Slice(snapshot, func(i, j int) bool {
			return snapshot[i].StartedAt.Before(snapshot[j].StartedAt)
		})

		dataMu.Lock()
		activeTasks = snapshot
		count := len(activeTasks)
		dataMu.Unlock()

		if count == 0 {
			listContainer.Hide()
			emptyState.Show()
		} else {
			emptyState.Hide()
			listContainer.Show()
			list.Refresh()
		}
	}

	// Initial Update
	updateList()

	// Event-Driven Refresher loop
	// Note: We still might want a slow ticker (e.g. 1s) just to update the "Duration" labels
	// because they are relative times. Pure event-driven only updates on add/remove.
	// But the requirement says "Remove Polling Ticker".
	// Compromise: Update on Event + Update on slow Ticker (1s) for liveness if list is not empty.
	// Or strictly follow "Trigger only when task added/removed".
	// The prompt said: "Event-Driven Radar (Remove Polling)".
	// I will remove the fast polling.

	go func() {
		// Slow ticker for duration updates only
		durationTicker := time.NewTicker(1 * time.Second)
		defer durationTicker.Stop()

		for {
			select {
			case <-state.Ctx.Done():
				return
			case <-state.UpdateChan:
				updateList()
			case <-durationTicker.C:
				// Refresh only if we have items, to update timers
				dataMu.Lock()
				count := len(activeTasks)
				dataMu.Unlock()
				if count > 0 {
					list.Refresh() // Just redraw content (Duration strings)
				}
			}
		}
	}()

	return container.NewPadded(mainStack)
}
