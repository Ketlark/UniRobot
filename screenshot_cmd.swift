import Foundation
import ScreenCaptureKit
import AppKit
import CoreImage
import CoreGraphics
import ImageIO
import UniformTypeIdentifiers

// Captures a window and returns its PNG data
func captureWindowArea(_ title: String) async -> Data? {
        CGMainDisplayID() // Initialize CoreGraphics
        
        do {
            // Retrieve the available windows
            let shareableContent = try await SCShareableContent.excludingDesktopWindows(false, onScreenWindowsOnly: true)
            
            // Find the first window that starts with the given title prefix
            guard let targetWindow = shareableContent.windows.first(where: { $0.title?.hasPrefix(title) == true }) else {
                print("No window found with the given prefix.")
                return nil
            }

            // Créer un filtre de capture pour la fenêtre choisie
            let contentFilter = SCContentFilter(desktopIndependentWindow: targetWindow)

            // Configuration du stream
            let streamConfig = SCStreamConfiguration()
            streamConfig.width = Int(targetWindow.frame.width)
            streamConfig.height = Int(targetWindow.frame.height)
            streamConfig.pixelFormat = kCVPixelFormatType_32BGRA  // Format d'image
            streamConfig.scalesToFit = true
            streamConfig.minimumFrameInterval = CMTime(value: 1, timescale: 10) // 1 FPS

            // Création du stream
            let stream = SCStream(filter: contentFilter, configuration: streamConfig, delegate: nil)
            try await stream.startCapture()

            // Capture d'une image
            if let screenshot = try? await SCScreenshotManager.captureImage(contentFilter: contentFilter, configuration: streamConfig) {
                print("Capture réussie !")

                if let croppedImage = cropCGImage(image: screenshot, to: CGRect(x: 8, y: 160, width: 128, height: 45)) {
                  // Enregistrer l'image sur le bureau
                  //saveImage(croppedImage)

                  //print("Image cropped successfully!")
                  return convertCGImageToPNG(croppedImage)
                }
            } else {
                print("Échec de la capture.")
            }

            // Arrêter le stream
            try await stream.stopCapture()
        } catch {
            print("Erreur : \(error.localizedDescription)")
        }

    return nil    
} 

// Converts a CGImage to PNG Data
func convertCGImageToPNG(_ image: CGImage) -> Data? {
    let mutableData = NSMutableData()
    guard let destination = CGImageDestinationCreateWithData(mutableData as CFMutableData, UTType.png.identifier as CFString, 1, nil) else {
        //print("Failed to create image destination for PNG conversion.")
        return nil
    }
    
    CGImageDestinationAddImage(destination, image, nil)
    
    // Finalize the image destination to write the data
    if CGImageDestinationFinalize(destination) {
        //print("Image successfully converted to PNG, size = \(mutableData.length) bytes")
        return mutableData as Data
    } else {
        //print("Failed to finalize PNG conversion.")
        return nil
    }
}

// Utility to crop a CGImage to a specific CGRect
func cropCGImage(image: CGImage, to rect: CGRect) -> CGImage? {
    return image.cropping(to: rect)
}

func saveImage(_ image: CGImage) {
    let desktopPath = FileManager.default.homeDirectoryForCurrentUser.appendingPathComponent("Desktop/screenshot.png")
        
    let destination = CGImageDestinationCreateWithURL(desktopPath as CFURL, "public.png" as CFString, 1, nil)!
    CGImageDestinationAddImage(destination, image, nil)
    CGImageDestinationFinalize(destination)
        
    print("Image enregistrée : \(desktopPath.path)")
}

if let imageData = await captureWindowArea("Ketlark") {
    // Output image data as bytes to standard output
    FileHandle.standardOutput.write(imageData)
}



