package capture

import (
	"fmt"
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

func CaptureOSXWindow() {
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
}
