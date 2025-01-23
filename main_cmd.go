package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Run the Swift script to capture the screenshot
	cmd := exec.Command("./screenshot")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Swift script:", err)
		return
	}

	// The output (image data) from the Swift script
	imageData := out.Bytes()
	fmt.Printf("Received %d bytes of image data\n", len(imageData))

	// Now you can save this data to a file (e.g., screenshot.png)
	// err = os.WriteFile("screenshot.png", imageData, 0644)
	// if err != nil {
	// 	fmt.Println("Error saving image:", err)
	// 	return
	// }

	fmt.Println("Image saved to screenshot.png")
}
