import AppKit
import SwiftUI

/// Colors the teaching views add to Theme.
enum TeachingStyle {
    /// A digit on a `.quaternary` circle or capsule, such as an answer's key or the review
    /// count: on that fill Theme.secondaryText reads only about 4 to 4.5:1, so this gray is a
    /// little deeper in light mode and lighter in dark mode.
    static let digitOnFill = Color(nsColor: NSColor(name: nil) { appearance in
        appearance.bestMatch(from: [.darkAqua, .aqua]) == .darkAqua
            ? NSColor(srgbRed: 0xC4 / 255, green: 0xC4 / 255, blue: 0xC4 / 255, alpha: 1)
            : NSColor(srgbRed: 0x4A / 255, green: 0x4A / 255, blue: 0x4A / 255, alpha: 1)
    })
}
