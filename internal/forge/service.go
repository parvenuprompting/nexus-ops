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
// - error: if immediate setup fails
func (p *Processor) Start(ctx context.Context, inputDir string, workerCount int) (<-chan int, <-chan struct{}, error) {
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
		go p.worker(ctx, &wg, jobs, outputDir, i)
	}

	// 5. Monitor Workers Completion
	go func() {
		wg.Wait()
		close(workersDone)
	}()

	return scanFound, workersDone, nil
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, outputDir string, workerID int) {
	defer wg.Done()

	// Register worker
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

			// Task Registration (Thread-Safe)
			taskID := fmt.Sprintf("img-%d-%s", workerID, filepath.Base(path))
			task := &sys.Task{
				ID:        taskID,
				Type:      sys.TaskTypeImageResize,
				StartedAt: time.Now(),
				Status:    "Processing",
			}
			p.State.AddActiveTask(task)

			err := p.processImageSafe(path, outputDir)

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
func (p *Processor) processImageSafe(srcPath, outDir string) error {
	// Acquire semaphore
	p.decodeSem <- struct{}{}
	defer func() { <-p.decodeSem }()

	// Open file
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode Config first (Memory Safety check - optional step but good practice)
	// For now we trust the semaphore to limit concurrent loading.

	// Decode
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Release semaphore early?
	// No, we hold it during resize too because resize creates a new buffer.
	// Actually, we could release after Decode if we want to allow resizing in parallel but limit decoding.
	// But `image.Decode` produces the big buffer. `draw.ApproxBiLinear.Scale` produces another one.
	// We'll be conservative and hold it.

	// Resize (50%)
	bounds := img.Bounds()
	newW := bounds.Dx() / 2
	newH := bounds.Dy() / 2
	rect := image.Rect(0, 0, newW, newH)

	dst := image.NewRGBA(rect)

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
