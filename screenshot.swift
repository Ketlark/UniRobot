import Foundation
import ScreenCaptureKit
import AppKit
import CoreImage
import CoreGraphics
import ImageIO
import UniformTypeIdentifiers

func cropCGImage(image: CGImage, to rect: CGRect) -> CGImage? {
    return image.cropping(to: rect)
}

@objcMembers
public class ScreenshotHelper: NSObject {
    static var shareableContent: SCShareableContent? = nil

    static var stream: SCStream? = nil
    static var streamConfig: SCStreamConfiguration? = nil

    static var contentFilter: SCContentFilter? = nil

    @objc public static func initCaptureWindow() {
        CGMainDisplayID()
        Task {
            do {
                // Initialize sharedContent once and keep it alive
                // Récupérer les fenêtres disponibles
                shareableContent = try await SCShareableContent.excludingDesktopWindows(false, onScreenWindowsOnly: true)
                print("Initialized shared content.")
            } catch {
                print("Failed to initialize shared content: \(error)")
            }

            // Trouver la première fenêtre qui commence par "Dofus"
            guard let targetWindow = shareableContent!.windows.first(where: { $0.title?.hasPrefix("Ketlark") == true }) else {
                print("Aucune fenêtre trouvée avec le préfixe")
                return
            }
                
            print("Fenêtre sélectionnée : \(targetWindow.title ?? "Sans titre")")



            // Créer un filtre de capture pour la fenêtre choisie
            contentFilter = SCContentFilter(desktopIndependentWindow: targetWindow)

            // Configuration du stream
            streamConfig = SCStreamConfiguration()
            streamConfig!.width = Int(targetWindow.frame.width * 2)
            streamConfig!.height = Int(targetWindow.frame.height * 2)
            streamConfig!.pixelFormat = kCVPixelFormatType_32BGRA  // Format d'image
            streamConfig!.scalesToFit = true
            streamConfig!.captureResolution = SCCaptureResolutionType.automatic
            streamConfig!.minimumFrameInterval = CMTime(value: 1, timescale: 1) // 1 FPS

            // Création du stream
            stream = SCStream(filter: contentFilter!, configuration: streamConfig!, delegate: nil)
            try await stream!.startCapture()
        }
    }

    @objc public static func captureWindowSync(sharedBuffer: UnsafeMutableRawPointer, bufferSize: UInt32) -> Bool {
        let resultBufferWrite = false
        let group = DispatchGroup()
        group.enter()

        Task {
            let data = await captureWindow(stream: stream)
            if let validData = data {
                let imageSize = min(validData.count, Int(bufferSize)) // Prevent overflow
                
                validData.copyBytes(to: sharedBuffer.assumingMemoryBound(to: UInt8.self), count: imageSize)
               // print("Swift: Captured Data Length =", validData.count)
            } else {
                print("Swift: Capture failed, got nil data")
            }
            group.leave()
        }

        group.wait()
        return !resultBufferWrite
    }

    @objc public static func captureWindow(stream: SCStream?) async -> Data? {
        do { 
            // Capture d'une image
            if let screenshot = try? await SCScreenshotManager.captureImage(contentFilter: contentFilter!, configuration: streamConfig!) {
                print("Capture réussie !")

                if let croppedImage = cropCGImage(image: screenshot, to: CGRect(x: 8, y: 160, width: 128, height: 45)) {
                  // Enregistrer l'image sur le bureau
                  // saveImage(croppedImage)

                  print("Image cropped successfully!")
                  return convertCGImageToPNG(croppedImage)
                }
            } else {
                print("Échec de la capture.")
            }

            // Arrêter le stream
            try await stream!.stopCapture()
        } catch {
            print("Erreur : \(error.localizedDescription)")
        }

        return nil
    }

    @objc private static func saveImage(_ image: CGImage) {
        let desktopPath = FileManager.default.homeDirectoryForCurrentUser.appendingPathComponent("Desktop/screenshot.png")
        
        let destination = CGImageDestinationCreateWithURL(desktopPath as CFURL, "public.png" as CFString, 1, nil)!
        CGImageDestinationAddImage(destination, image, nil)
        CGImageDestinationFinalize(destination)
        
        print("Image enregistrée : \(desktopPath.path)")
    }

    private static func convertCGImageToPNG(_ image: CGImage) -> Data? {
        let mutableData = NSMutableData()
        guard let destination = CGImageDestinationCreateWithData(mutableData as CFMutableData, UTType.png.identifier as CFString, 1, nil) else {
            print("Swift: Failed to create image destination")
            return nil
        }

        CGImageDestinationAddImage(destination, image, nil)
        if CGImageDestinationFinalize(destination) {
            print("Swift: Image converted to PNG, size =", mutableData.length)
            return mutableData as Data
        } else {
            print("Swift: Failed to finalize image destination")
            return nil
        }
    }
}
