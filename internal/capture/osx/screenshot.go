package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	memory "unirobot/util"

	"github.com/ebitengine/purego"
)

// const SWIFT_FREE_MEMORY_FUNC_SYMBOL = "$s16ScreenshotHelperAAC16freeSharedMemoryyySvFZ"
const SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL = "$s16ScreenshotHelperAAC17captureWindowSync12sharedBuffer10bufferSizeSbSv_s6UInt32VtFZ"
const SWIFT_INIT_CAPTURE_FUNC_SYMBOL = "$s16ScreenshotHelperAAC17initCaptureWindowyyFZ"

const MAP_POS_BUFFER_SIZE = 50 * 1024 //50KB

func CaptureOSXWindow(bufferShared []byte, libc uintptr) {
	// Find the Swift function that returns NSData*
	intptr, err := purego.Dlsym(libc, SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL)
	if err != nil || intptr == 0 {
		fmt.Println("❌ Failed to find symbol captureWindowSync")
		return
	}
	fmt.Println("✅ Symbol found successfully")

	// Load Swift function returning pointer to shared buffer memory
	var captureWindow func(sharedBuffer unsafe.Pointer, bufferSize uint32) bool
	purego.RegisterLibFunc(&captureWindow, libc, SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL)

	// Call Swift function
	result := captureWindow(unsafe.Pointer(&bufferShared[0]), uint32(len(bufferShared)))
	if !result {
		fmt.Println("❌ No screenshot data received.")
		return
	}
}

func main() {
	bufferShared := memory.CreateSharedMemory(MAP_POS_BUFFER_SIZE)
	libc, err := purego.Dlopen("./ScreenshotHelper.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	//Init capture window
	intptr, err := purego.Dlsym(libc, SWIFT_INIT_CAPTURE_FUNC_SYMBOL)
	if err != nil || intptr == 0 {
		fmt.Println("❌ Failed to find symbol captureWindowSync")
		return
	}
	//fmt.Println("✅ Symbol found successfully")

	// Load Swift function returning pointer to shared buffer memory
	var initCaptureWindow func()
	purego.RegisterLibFunc(&initCaptureWindow, libc, SWIFT_INIT_CAPTURE_FUNC_SYMBOL)
	initCaptureWindow()

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			CaptureOSXWindow(bufferShared, libc)
			//fmt.Println("Screenshot taken: ", time.Now())
		}
	}()

	// Handle Ctrl+C (SIGINT) to free memory before exiting
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal
	go func() {
		<-sigChan
		fmt.Println("\nSignal received, freeing memory...")
		memory.FreeSharedMemory(bufferShared)
		os.Exit(0)
	}()

	for {
		time.Sleep(2000 * time.Millisecond)
		// Save to file
		// err := os.WriteFile("/Users/dehoux/Desktop/screenshot.png", bufferShared, 0644)
		// if err != nil {
		// 	fmt.Println("❌ Error saving image:", err)
		// 	return
		// }

		fmt.Println("✅ Image saved successfully! Size:", "bytes")
	}
}
