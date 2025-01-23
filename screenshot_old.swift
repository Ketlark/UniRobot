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

let BUFFER_SIZE = 50 * 1024 // 50KB buffer

func createSharedMemory() -> UnsafeMutableRawPointer? {
    let memory = mmap(nil, BUFFER_SIZE, PROT_READ | PROT_WRITE, MAP_SHARED | MAP_ANONYMOUS, -1, 0)
    if memory == MAP_FAILED {
        print("mmap failed")
        return nil
    }

    print("Swift: Shared memory created at \(memory!)")

    // Example: Fill memory with some dummy data
    //memset(memory, 255, BUFFER_SIZE) // Fill with white pixels (RGBA: 255,255,255,255)
    
    return memory
}

@objcMembers
public class ScreenshotHelper: NSObject {
    @objc public static func captureWindowSync() -> UnsafeMutableRawPointer? {
        var resultSharedPointerAddress: UnsafeMutableRawPointer?
        let group = DispatchGroup()
        group.enter()

        Task {
            let data = await captureWindow()
            if let validData = data {
                let imageSize = min(validData.count, BUFFER_SIZE) // Prevent overflow
                
                if let sharedBuffer = createSharedMemory() {
                    validData.copyBytes(to: sharedBuffer.assumingMemoryBound(to: UInt8.self), count: imageSize)

                    print("Swift: Captured Data Length =", validData.count)
                    resultSharedPointerAddress = sharedBuffer
                }
            } else {
                print("Swift: Capture failed, got nil data")
            }
            group.leave()
        }

        group.wait()
        return resultSharedPointerAddress
    }

    @objc public static func captureWindow() async -> Data? {
        do {
            CGMainDisplayID()
                
            // Récupérer les fenêtres disponibles
            let shareableContent = try await SCShareableContent.excludingDesktopWindows(false, onScreenWindowsOnly: true)
                
            // Afficher toutes les fenêtres (Debug)
            // print("Fenêtres disponibles :")
            // for window in shareableContent.windows {
            //     print("ID: \(window.windowID), Titre: \(window.title ?? "Sans titre")")
            // }
                
            // Trouver la première fenêtre qui commence par "Dofus"
            guard let targetWindow = shareableContent.windows.first(where: { $0.title?.hasPrefix("Ketlark") == true }) else {
                print("Aucune fenêtre trouvée avec le préfixe")
                return nil
            }
                
            print("Fenêtre sélectionnée : \(targetWindow.title ?? "Sans titre")")

            // Créer un filtre de capture pour la fenêtre choisie
            let contentFilter = SCContentFilter(desktopIndependentWindow: targetWindow)

            // Configuration du stream
            let streamConfig = SCStreamConfiguration()
            streamConfig.width = Int(targetWindow.frame.width * 2)
            streamConfig.height = Int(targetWindow.frame.height * 2)
            streamConfig.pixelFormat = kCVPixelFormatType_32BGRA  // Format d'image
            streamConfig.scalesToFit = true
            streamConfig.captureResolution = SCCaptureResolutionType.automatic
            streamConfig.minimumFrameInterval = CMTime(value: 1, timescale: 1) // 1 FPS

            // Création du stream
            let stream = SCStream(filter: contentFilter, configuration: streamConfig, delegate: nil)
            try await stream.startCapture()
            
            // Capture d'une image
            if let screenshot = try? await SCScreenshotManager.captureImage(contentFilter: contentFilter, configuration: streamConfig) {
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
            try await stream.stopCapture()
        } catch {
            print("Erreur : \(error.localizedDescription)")
        }

        return nil
    }

    @objc public static func freeSharedMemory(_ pointer: UnsafeMutableRawPointer) {
        let result = munmap(pointer, BUFFER_SIZE) // Only if mmap was used
        if result != 0 {
            print("munmap failed")
            return
        }
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
