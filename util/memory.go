package memory

import (
	"fmt"
	"syscall"
)

func CreateSharedMemory(bufferSize uint32) []byte {
	shm, err := syscall.Mmap(-1, 0, int(bufferSize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("Erreur allocation mémoire partagée:", err)
		return nil
	}

	return shm
}

func FreeSharedMemory(ptr []byte) {
	err := syscall.Munmap(ptr)
	if err != nil {
		fmt.Println("Erreur libération mémoire partagée:", err)
		return
	}
}
