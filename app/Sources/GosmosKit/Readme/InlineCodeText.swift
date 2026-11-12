import Foundation

/// Inline code as a readme shows it. A call, a path or a command broken at a line's end reads as
/// two pieces of code, so a span short enough to fit a line keeps to one: a word joiner between
/// each two of its characters leaves the line nowhere to break inside it. A longer span still
/// breaks where it must. Copied code carries the joiners along, so a paste into the editor
/// leaves them out (CodeEditor's PastedText), and so does leaving the app (AppDelegate).
public enum InlineCodeText {
    /// The longest span kept to one line, in characters: well within what the narrowest panel
    /// holds.
    public static let unbrokenLimit = 40

    /// `code` as shown.
    public static func shown(_ code: String) -> String {
        guard code.count > 1, code.count <= unbrokenLimit else { return code }
        return code.map(String.init).joined(separator: "\u{2060}")
    }
}
