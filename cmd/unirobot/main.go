package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unirobot/internal/capture"
)

func main() {
	cm := capture.NewCaptureManager()

	println("Good cm")

	go cm.ProcessCapture()

	// Handle Ctrl+C (SIGINT) to free memory before exiting
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal
	go func() {
		<-sigChan
		fmt.Println("\nSignal received, freeing memory...")
		cm.Free()
		os.Exit(0)
	}()

	for {
		time.Sleep(5000 * time.Millisecond)
		// Save to file
		// err := os.WriteFile("/Users/dehoux/Desktop/screenshot.png", bufferShared, 0644)
		// if err != nil {
		// 	fmt.Println("❌ Error saving image:", err)
		// 	return
		// }

		_, size := cm.GetGlobalBuffer()

		fmt.Println("✅ Image saved successfully! Size:", size)
	}
}
