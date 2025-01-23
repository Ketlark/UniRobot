package capture

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/fstanis/screenresolution"
	"github.com/go-vgo/robotgo"
)

func TakeScreenshot(options ...uint) bytes.Buffer {
	// Créer un buffer pour capturer l'image
	var buf bytes.Buffer

	// Créer un buffer pour capturer les erreurs
	var stderr bytes.Buffer

	//Recuperer la resolution de l'écran principal
	resolution := screenresolution.GetPrimary()
	if resolution == (&screenresolution.Resolution{}) {
		log.Fatalf("failed to get screen resolution")
		resolution = &screenresolution.Resolution{Width: 1920, Height: 1080} //Based on classic monitor resolution
	}
	// Obtenez l'ID de la fenêtre active
	window := robotgo.GetScreenSize()

	x, y, w, h := 0, resolution.Height, 800, 600
	cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d", w, h, x, y)

	robotgo
	// Lancer FFmpeg pour capturer une image de l'écran et rediriger la sortie vers le buffer
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation", // Utiliser le framework AVFoundation pour macOS
		"-i", "3", // Capture l'écran principal (ajuste ce numéro si nécessaire)
		"-vframes", "1", // Capture une seule image
		"-video_size", "1920x1080", // Taille de la capture
		"-pix_fmt", "rgb24", // Utilisation d'un format RGB compatible
		"-vcodec", "png",
		"-vf", cropFilter, // Dynamic cropping
		"-f", "image2pipe", // Sortie au format image2pipe
		"pipe:1", // Envoie la sortie directement dans un flux
	)

	// Rediriger la sortie de la commande vers le buffer
	cmd.Stdout = &buf
	cmd.Stderr = &stderr

	// Exécution de la commande FFmpeg
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Erreur lors de la capture d'image : %v\nstderr: %s", err, stderr.String())
	}

	return buf
}
