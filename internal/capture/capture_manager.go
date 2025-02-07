package capture

import (
	"runtime"
	"time"
	capture_mac "unirobot/internal/capture/osx"
	"unirobot/internal/config"
	"unirobot/pkg/utils"
)

type CaptureManager struct {
	globalBuffer *utils.BufferManager
	areas        map[string]*utils.BufferManager
}

func NewCaptureManager() *CaptureManager {
	cm := &CaptureManager{
		globalBuffer: utils.NewBufferManager(config.GLOBAL_BUFFER_SIZE),
		areas:        make(map[string]*utils.BufferManager),
	}

	// Initialize areas from config
	for _, area := range config.GetGameAreas() {
		cm.areas[area.Name] = utils.NewBufferManager(area.Size)
	}

	if runtime.GOOS == "darwin" {
		capture_mac.InitLibrary()
	}

	return cm
}

/*
 * This method take a screenshot and put it in the global buffer. Based on OS, screenshot will be processed in specific method to use GPU power to treat the image.
 */
func (cm *CaptureManager) ProcessCapture() {
	for {
		time.Sleep(1000 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			//cm.processWindowsCapture()
		case "darwin":
			capture_mac.TakeScreenshot(cm.globalBuffer.Buffer)
		case "linux":
			//cm.processLinuxCapture()
		}
	}
}

func (cm *CaptureManager) GetGlobalBuffer() ([]byte, uint32) {
	return cm.globalBuffer.GetBuffer()
}

func (cm *CaptureManager) GetAreaBuffer(areaName string) ([]byte, uint32) {
	if area, ok := cm.areas[areaName]; ok {
		return area.GetBuffer()
	}
	return nil, 0
}

func (cm *CaptureManager) Free() {
	cm.globalBuffer.Free()
	for _, area := range cm.areas {
		area.Free()
	}
}
