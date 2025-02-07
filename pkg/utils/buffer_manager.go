package utils

type BufferManager struct {
	Buffer []byte
	Size   uint32
}

func NewBufferManager(size uint32) *BufferManager {
	return &BufferManager{
		Buffer: CreateSharedMemory(size),
		Size:   size,
	}
}

func (bm *BufferManager) GetBuffer() ([]byte, uint32) {
	return bm.Buffer, bm.Size
}

func (bm *BufferManager) Free() {
	FreeSharedMemory(bm.Buffer)
}
