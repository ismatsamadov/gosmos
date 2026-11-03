import AppKit
import GosmosKit
import SwiftUI

/// Colors and type sizes from docs/design.md.
enum Theme {
    /// Go Teal #00758D in light mode, Go Blue #00ADD8 in dark mode. (An asset catalog's
    /// AccentColor needs actool, which ships with Xcode only, so the app sets it in code.)
    static let accentNS = dynamic(light: 0x00758D, dark: 0x00ADD8, name: "GoAccent")
    static let accent = Color(nsColor: accentNS)

    /// The fill of the prominent action. White reads at only 2.5:1 on Go Blue, so in dark mode
    /// the fill is a deeper Go blue that keeps white labels above 4.5:1; light mode keeps Go Teal.
    static let prominent = Color(nsColor: dynamic(light: 0x00758D, dark: 0x007093, name: "GoProminent"))

    // Text colors for results, warnings and callouts: the system hues, deepened in light mode
    // and lightened in dark mode so text keeps 4.5:1 on window, panel and glass backgrounds.
    // The light green is deep enough (7.8:1 on white) that thin, small text still measures
    // above 5:1 as drawn.
    static let successNS = dynamic(light: 0x0A5F2C, dark: 0x32D74B, name: "GoSuccess")
    static let failureNS = dynamic(light: 0xB42318, dark: 0xFF7A70, name: "GoFailure")
    static let success = Color(nsColor: successNS)
    static let failure = Color(nsColor: failureNS)
    static let warning = Color(nsColor: dynamic(light: 0x8A5300, dark: 0xFFB340, name: "GoWarning"))
    static let note = Color(nsColor: dynamic(light: 0x0B5CBF, dark: 0x5AA9FF, name: "GoNote"))
    static let important = Color(nsColor: dynamic(light: 0x7E3BA6, dark: 0xD88DF7, name: "GoImportant"))
    /// Secondary text. The system's secondary label reads about 3.9:1 on white in light mode;
    /// this one keeps 4.5:1 on window, panel and sidebar backgrounds, and on a `.quaternary`
    /// capsule or a gray selection over them, as the neutral gray it replaced did: a cool gray
    /// of the same lightness (#606B78 was lighter and fell to about 4:1 there). Symbols keep
    /// `.secondary`, which is enough for them at 3:1.
    static let secondaryText = Color(nsColor: dynamic(light: 0x56606C, dark: 0xA7B1BD, name: "GoSecondaryText"))
    /// Words a search matched, in the sidebar's results. On the sidebar's gray Go Teal reads
    /// only 4.6:1, so light mode deepens it (6.3:1); dark mode keeps Go Blue.
    static let searchMark = Color(nsColor: dynamic(light: 0x005E70, dark: 0x00ADD8, name: "GoSearchMark"))

    static let canvas = Color(nsColor: dynamic(light: 0xF5F7F9, dark: 0x15191F, name: "MissionCanvas"))
    static let surface = Color(nsColor: dynamic(light: 0xFFFFFF, dark: 0x1D232C, name: "MissionSurface"))
    static let inset = Color(nsColor: dynamic(light: 0xEEF2F5, dark: 0x181D24, name: "MissionInset"))
    static let sidebar = Color(nsColor: dynamic(light: 0xEEF1F4, dark: 0x1A1F26, name: "MissionSidebar"))
    static let border = Color(nsColor: dynamic(light: 0xDFE5EB, dark: 0x333C48, name: "MissionBorder"))
    /// The outline of a card, an answer or a pane: the fine border, or with Increase Contrast one
    /// that reads against the canvas (the border alone measures about 1.2:1 there).
    static func outline(_ contrast: ColorSchemeContrast) -> Color {
        contrast == .increased ? secondaryText : border
    }
    /// The welcome panel has one deliberately dark surface in either appearance.
    static let ink = Color(red: 23 / 255, green: 43 / 255, blue: 54 / 255)
    static let onInk = Color(red: 182 / 255, green: 204 / 255, blue: 212 / 255)
    static let orbit = Color(red: 103 / 255, green: 218 / 255, blue: 227 / 255)
    static let gold = Color(nsColor: dynamic(light: 0x926100, dark: 0xF4C66A, name: "MissionGold"))
    static let violet = Color(nsColor: dynamic(light: 0x7051B4, dark: 0xB9A0F6, name: "MissionViolet"))

    static let bodySize: CGFloat = 15
    static let codeSize: CGFloat = 13
    /// The reading column: room for 80 columns of code.
    static let readableWidth: CGFloat = 680
    /// The measure of body text in the column, about 80 characters of SF Pro at 15 pt; code,
    /// tables and panels take the whole column.
    static let proseWidth: CGFloat = 560

    static func headingSize(_ level: Int) -> CGFloat {
        switch level {
        case 1: 28
        case 2: 22
        case 3: 18
        default: 15
        }
    }

    /// The callout colors; info and other kinds look like a note, danger like a caution.
    static func calloutColor(_ kind: String) -> Color {
        switch Readme.calloutKind(kind) {
        case "tip": success
        case "important": important
        case "warning": warning
        case "caution": failure
        default: note
        }
    }

    /// An sRGB color for light mode and another for dark mode.
    private static func dynamic(light: UInt32, dark: UInt32, name: String) -> NSColor {
        func color(_ value: UInt32) -> NSColor {
            NSColor(srgbRed: CGFloat((value >> 16) & 0xFF) / 255, green: CGFloat((value >> 8) & 0xFF) / 255,
                    blue: CGFloat(value & 0xFF) / 255, alpha: 1)
        }
        let lightColor = color(light), darkColor = color(dark)
        return NSColor(name: name) { appearance in
            appearance.bestMatch(from: [.darkAqua, .aqua]) == .darkAqua ? darkColor : lightColor
        }
    }
}

extension Duration {
    var milliseconds: Double {
        Double(components.seconds) * 1000 + Double(components.attoseconds) / 1e15
    }

    /// "375 ms" or "2.4 s".
    var shortText: String {
        milliseconds < 1000 ? String(format: "%.0f ms", milliseconds) : String(format: "%.1f s", milliseconds / 1000)
    }
}

extension View {
    /// A secondary action in content: plain glass. Under the window's accent tint a .glass
    /// button fills with the accent like the prominent action; a nil tint gives back the
    /// untinted glass in light and dark mode.
    func secondaryGlassButton() -> some View {
        buttonStyle(.glass).tint(nil)
    }

    /// The prominent action: glass filled with `Theme.prominent`, white label.
    func prominentGlassButton() -> some View {
        buttonStyle(.glassProminent).tint(Theme.prominent)
    }

    /// Shows the view only once it has been on screen for a moment, so a spinner for a wait
    /// that ends quickly never flashes.
    func appearsAfterDelay(_ delay: Duration = .milliseconds(400)) -> some View {
        modifier(AppearsAfterDelay(delay: delay))
    }

    /// Turns the spinner symbol (progress.indicator) while work runs; with Reduce Motion its
    /// spokes light in turn instead. Used in place of ProgressView's spinner, which the
    /// system draws dimmed in a window that is not active (about 2:1 in light mode).
    func runningSymbolEffect(_ isActive: Bool) -> some View {
        modifier(RunningSymbolEffect(isActive: isActive))
    }
}

private struct RunningSymbolEffect: ViewModifier {
    let isActive: Bool
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    func body(content: Content) -> some View {
        content
            .symbolEffect(.rotate.clockwise, options: .repeat(.continuous), isActive: isActive && !reduceMotion)
            .symbolEffect(.variableColor.iterative.nonReversing, options: .repeating, isActive: isActive && reduceMotion)
    }
}

private struct AppearsAfterDelay: ViewModifier {
    let delay: Duration
    @ViewState private var shown = false

    func body(content: Content) -> some View {
        content
            .opacity(shown ? 1 : 0)
            .task {
                try? await Task.sleep(for: delay)
                shown = true
            }
    }
}
