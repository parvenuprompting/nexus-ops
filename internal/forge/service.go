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
	// Phase 3: Memory Safety - Semaphore
	decodeSem chan struct{}
}

// NewProcessor creates a new Processor.
func NewProcessor(state *sys.AppState) *Processor {
	// Allow max 4 concurrent heavy decodes
	return &Processor{
		State:     state,
		decodeSem: make(chan struct{}, 4),
	}
}

// Start initiates the worker pool and streaming scanner.
// Returns:
// - scanFound: channel that receives total count when scanning is done
// - workersDone: channel that closes when all workers are finished
// - workersDone: channel that closes when all workers are finished
// - error: if immediate setup fails
func (p *Processor) Start(ctx context.Context, inputDir string, workerCount int, settings sys.ForgeSettings) (<-chan int, <-chan struct{}, error) {
	// 1. Prepare Output Directory
	outputDir := filepath.Join(inputDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create output dir: %w", err)
	}

	// 2. Setup Channels
	// Unbuffered or small buffer for streaming jobs
	jobs := make(chan string, workerCount*2)
	scanFound := make(chan int, 1)
	workersDone := make(chan struct{})

	// 3. Start Streaming Scanner
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC in Scanner: %v\n", r)
			}
		}()
		defer close(jobs)
		total, err := scanImagesStreaming(ctx, inputDir, jobs)
		if err != nil {
			fmt.Printf("Scanner Error: %v\n", err) // Log for now
		}
		scanFound <- total
		close(scanFound)
	}()

	// 4. Start Workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, jobs, outputDir, i, settings)
	}

	// 5. Monitor Workers Completion
	go func() {
		wg.Wait()
		close(workersDone)
	}()

	return scanFound, workersDone, nil
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, outputDir string, workerID int, settings sys.ForgeSettings) {
	defer wg.Done()

	// Register worker
	p.State.Metrics.ActiveGoroutines.Add(1)

	defer func() {
		p.State.Metrics.ActiveGoroutines.Add(-1)
		if r := recover(); r != nil {
			fmt.Printf("PANIC in Worker %d: %v\n", workerID, r)
			p.State.Metrics.Errors.Add(1)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-jobs:
			if !ok {
				return
			}

			// Task Registration (Thread-Safe)
			taskID := fmt.Sprintf("img-%d-%s", workerID, filepath.Base(path))
			task := &sys.Task{
				ID:        taskID,
				Name:      filepath.Base(path),
				Type:      sys.TaskTypeImageResize,
				StartedAt: time.Now(),
				Status:    "Processing",
			}
			p.State.AddActiveTask(task)

			err := p.processImageSafe(path, outputDir, settings)

			if err != nil {
				p.State.Metrics.Errors.Add(1)
			} else {
				p.State.Metrics.TasksCompleted.Add(1)
			}

			p.State.RemoveActiveTask(taskID)
		}
	}
}

// scanImagesStreaming walks the directory and sends files to the jobs channel
func scanImagesStreaming(ctx context.Context, dir string, jobs chan<- string) (int, error) {
	count := 0
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, e := range entries {
		if ctx.Err() != nil {
			return count, ctx.Err()
		}
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			select {
			case jobs <- filepath.Join(dir, e.Name()):
				count++
			case <-ctx.Done():
				return count, ctx.Err()
			}
		}
	}
	return count, nil
}

// processImageSafe with Memory Semaphore
func (p *Processor) processImageSafe(srcPath, outDir string, settings sys.ForgeSettings) error {
	// Acquire semaphore
	p.decodeSem <- struct{}{}
	defer func() { <-p.decodeSem }()

	// Open file
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Calculate Target Dimensions (Crop vs Resize)
	bounds := img.Bounds()
	curW, curH := bounds.Dx(), bounds.Dy()

	// Default: Keep original resolution if no resize logic is specified?
	// The prompt implies "resizing" but focus is on Aspect Ratio.
	// We will implement CENTER CROP to the aspect ratio, maintaining MAX possible size.

	targetW, targetH := curW, curH

	if settings.AspectRatio != "" && settings.AspectRatio != "original" {
		parts := strings.Split(settings.AspectRatio, ":")
		if len(parts) == 2 {
			var wRatio, hRatio float64
			fmt.Sscanf(parts[0], "%f", &wRatio)
			fmt.Sscanf(parts[1], "%f", &hRatio)

			if wRatio > 0 && hRatio > 0 {
				ratio := wRatio / hRatio
				imgRatio := float64(curW) / float64(curH)

				if imgRatio > ratio {
					// Image is wider than target: Crop Width
					targetW = int(float64(curH) * ratio)
				} else {
					// Image is taller than target: Crop Height
					targetH = int(float64(curW) / ratio)
				}
			}
		}
	}

	// Center Crop Logic
	x0 := (curW - targetW) / 2
	y0 := (curH - targetH) / 2
	x1 := x0 + targetW
	y1 := y0 + targetH

	cropRect := image.Rect(x0, y0, x1, y1)

	// Create destination image
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	// Draw cropped area
	draw.Draw(dst, dst.Bounds(), img, cropRect.Min, draw.Src)

	// Save
	outName := filepath.Base(srcPath)
	// Handle format change
	ext := filepath.Ext(outName)
	nameWithoutExt := strings.TrimSuffix(outName, ext)

	saveFormat := settings.Format
	if saveFormat == "" {
		saveFormat = "jpg" // Default
	}

	outPath := filepath.Join(outDir, nameWithoutExt+"."+saveFormat)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch strings.ToLower(saveFormat) {
	case "jpeg", "jpg":
		return jpeg.Encode(outFile, dst, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(outFile, dst)
	default: // Default to jpg
		return jpeg.Encode(outFile, dst, &jpeg.Options{Quality: 90})
	}
}
