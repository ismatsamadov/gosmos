import Foundation

/// Characters text can carry unseen that never belong in code: a lesson's readme puts word
/// joiners inside inline code, to keep each span on one line, and narrow no-break spaces around
/// it; text from elsewhere can hold zero-width spaces or a byte order mark. Go refuses each of
/// them outside a string, and nobody can see which one is there.
public enum PastedText {
    static let invisible: Set<Unicode.Scalar> = ["\u{2060}", "\u{200B}", "\u{FEFF}", "\u{202F}"]

    /// Whether `text` holds any of them.
    public static func needsCleaning(_ text: String) -> Bool {
        text.unicodeScalars.contains { invisible.contains($0) }
    }

    /// `text` without them.
    public static func clean(_ text: String) -> String {
        var scalars = String.UnicodeScalarView()
        scalars.append(contentsOf: text.unicodeScalars.lazy.filter { !invisible.contains($0) })
        return String(scalars)
    }
}
