package main

import (
	"context"
	"fmt"
	"nexus-ops/internal/sys"
	"os"
	"testing"
	"time"
)

func TestCrashReproduction(t *testing.T) {
	// Setup
	ctx, cancel := context.WithCancel(context.Background())
	metrics := &sys.Metrics{}
	state := sys.NewAppState(ctx, cancel, metrics)
	app := NewApp(state)
	app.startup(context.Background())

	// Create dummy dir
	err := os.MkdirAll("test_images", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("test_images")

	// Create a dummy image file to trigger processing
	f, _ := os.Create("test_images/test.jpg")
	f.Close()

	// Trigger processing
	fmt.Println("Triggering StartProcessing...")
	app.StartProcessing("test_images")

	// Wait to see if it panics
	time.Sleep(2 * time.Second)
	fmt.Println("Survived StartProcessing")
}
