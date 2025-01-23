package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/ebitengine/purego"
)

const SWIFT_FREE_MEMORY_FUNC_SYMBOL = "$s16ScreenshotHelperAAC16freeSharedMemoryyySvFZ"
const SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL = "$s16ScreenshotHelperAAC17captureWindowSyncSvSgyFZ"

const BUFFER_SIZE = 50 * 1024 // 50KB buffer

func freeMemorySwift(libc uintptr, ptrToFree unsafe.Pointer) {
	// Find the Swift function
	intptr, err := purego.Dlsym(libc, SWIFT_FREE_MEMORY_FUNC_SYMBOL)
	if err != nil || intptr == 0 {
		fmt.Println("❌ Failed to find symbol captureWindowSync")
		return
	}
	fmt.Println("✅ Symbol found successfully")

	// Load Swift function returning NSData*
	var freeMemory func(ptr unsafe.Pointer)
	purego.RegisterLibFunc(&freeMemory, libc, SWIFT_FREE_MEMORY_FUNC_SYMBOL)

	// Call Swift function
	freeMemory(ptrToFree)
}

func main() {
	libc, err := purego.Dlopen("./ScreenshotHelper.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	// Find the Swift function that returns NSData*
	intptr, err := purego.Dlsym(libc, SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL)
	if err != nil || intptr == 0 {
		fmt.Println("❌ Failed to find symbol captureWindowSync")
		return
	}
	fmt.Println("✅ Symbol found successfully")

	// Load Swift function returning NSData*
	var captureWindow func() unsafe.Pointer
	purego.RegisterLibFunc(&captureWindow, libc, SWIFT_CAPTURE_WINDOW_FUNC_SYMBOL)

	// Call Swift function
	sharedBufferPtr := captureWindow()
	if sharedBufferPtr == nil {
		fmt.Println("❌ No screenshot data received.")
		return
	}

	// Convert the pointer to a byte slice
	imageData := unsafe.Slice((*byte)(sharedBufferPtr), BUFFER_SIZE) // Adjust BUFFER_SIZE as needed
	if len(imageData) == 0 {
		fmt.Println("❌ Buffer is empty.")
		return
	}

	// Print the pointer address in Go
	fmt.Println("Pointer Address in Go: ", sharedBufferPtr)

	// Save to file
	err = os.WriteFile("screenshot.png", imageData, 0644)
	if err != nil {
		fmt.Println("❌ Error saving image:", err)
		return
	}

	fmt.Println("✅ Image saved successfully! Size:", "bytes")

	freeMemorySwift(libc, sharedBufferPtr)

	fmt.Println("✅ Memory unmapped successfully")
}
