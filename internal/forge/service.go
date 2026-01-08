package forge

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"nexus-ops/internal/sys"

	"golang.org/x/image/draw"
)

// Processor handles the image resizing workload.
type Processor struct {
	State *sys.AppState
}

// NewProcessor creates a new Processor.
func NewProcessor(state *sys.AppState) *Processor {
	return &Processor{
		State: state,
	}
}

// Start initiates the worker pool to process images in inputDir.
// It returns a total count of images found, and an error if any setup fails.
func (p *Processor) Start(ctx context.Context, inputDir string, workerCount int) (int, error) {
	// 1. Validate Input
	matches, err := scanImages(inputDir)
	if err != nil {
		return 0, err
	}
	totalImages := len(matches)
	if totalImages == 0 {
		return 0, nil
	}

	// 2. Prepare Output Directory
	outputDir := filepath.Join(inputDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create output dir: %w", err)
	}

	// 3. Setup Channels and WaitGroup
	jobs := make(chan string, totalImages)
	var wg sync.WaitGroup

	// 4. Start Workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, jobs, outputDir, i)
	}

	// 5. Enqueue Jobs
	// We do this in a separate goroutine or just push them all since buffer is large enough.
	// Since we know the count, we buffered the channel to hold all.
	for _, path := range matches {
		jobs <- path
	}
	close(jobs)

	// 6. Wait for completion in background to cleanup
	// We don't block Start(), but we can return.
	// Actually, for the UI progress bar, we might want to know when it's *all* done?
	// The prompt implies we just fire it off. The UI updates via Metrics.

	// Wait logic is handled by the workers decrementing waitgroup?
	// To cleanly 'stop' the task group, we rely on context.

	return totalImages, nil
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, outputDir string, workerID int) {
	defer wg.Done()

	// Register worker as an active goroutine
	p.State.Metrics.ActiveGoroutines.Add(1)
	defer p.State.Metrics.ActiveGoroutines.Add(-1)

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-jobs:
			if !ok {
				return
			}

			// Register specific task in Radar
			taskID := fmt.Sprintf("img-%d-%s", workerID, filepath.Base(path))
			task := &sys.Task{
				ID:        taskID,
				Type:      sys.TaskType("Image-Resize"), // Cast string to TaskType
				StartedAt: time.Now(),
				Status:    "Processing",
			}
			p.State.ActiveTasks.Store(taskID, task)

			err := processImage(path, outputDir)

			// Update Metrics
			if err != nil {
				p.State.Metrics.Errors.Add(1)
				// fmt.Printf("Worker %d error on %s: %v\n", workerID, path, err)
			} else {
				p.State.Metrics.TasksCompleted.Add(1)
			}

			// Deregister from Radar
			p.State.ActiveTasks.Delete(taskID)
		}
	}
}

func scanImages(dir string) ([]string, error) {
	var images []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			images = append(images, filepath.Join(dir, e.Name()))
		}
	}
	return images, nil
}

func processImage(srcPath, outDir string) error {
	// Open file
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize (50%)
	bounds := img.Bounds()
	newW := bounds.Dx() / 2
	newH := bounds.Dy() / 2
	rect := image.Rect(0, 0, newW, newH)

	dst := image.NewRGBA(rect)

	// Draw with ApproxBiLinear (decent balance of speed/quality)
	draw.ApproxBiLinear.Scale(dst, rect, img, bounds, draw.Over, nil)

	// Save
	outName := filepath.Base(srcPath)
	outPath := filepath.Join(outDir, outName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(outFile, dst, &jpeg.Options{Quality: 80})
	case "png":
		return png.Encode(outFile, dst)
	default:
		return fmt.Errorf("unsupported format for saving: %s", format)
	}
}
