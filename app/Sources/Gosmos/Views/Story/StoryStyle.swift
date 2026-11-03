import GosmosKit
import SwiftUI

/// The story's colors: the night sky of the opening and the mission patches. They stay the same
/// in either appearance, as the welcome panel's ink does.
enum StoryPalette {
    static let night = Color(red: 11 / 255, green: 20 / 255, blue: 27 / 255)
    /// A card on the night sky.
    static let panel = Color(red: 19 / 255, green: 35 / 255, blue: 45 / 255)
    static let earth = Color(red: 14 / 255, green: 40 / 255, blue: 50 / 255)
    /// Text on a patch's color.
    static let patchInk = Color(red: 12 / 255, green: 24 / 255, blue: 31 / 255)

    /// A patch's color, by its chapter's number: six that read apart and keep dark text above
    /// 4.5:1, in turn.
    static func patch(_ number: Int) -> Color {
        let colors: [Color] = [
            Color(red: 103 / 255, green: 218 / 255, blue: 227 / 255), // orbit teal
            Color(red: 244 / 255, green: 198 / 255, blue: 106 / 255), // gold
            Color(red: 185 / 255, green: 160 / 255, blue: 246 / 255), // violet
            Color(red: 255 / 255, green: 150 / 255, blue: 128 / 255), // coral
            Color(red: 120 / 255, green: 214 / 255, blue: 140 / 255), // green
            Color(red: 122 / 255, green: 180 / 255, blue: 255 / 255), // sky
        ]
        return colors[(max(1, number) - 1) % colors.count]
    }
}

/// A story text, its paragraphs as body text with their inline Markdown.
struct StoryText: View {
    let text: String
    var size: CGFloat = 15
    var color: Color = .primary
    var spacing: CGFloat = 10

    var body: some View {
        VStack(alignment: .leading, spacing: spacing) {
            ForEach(Array(paragraphs.enumerated()), id: \.offset) { _, inlines in
                Text(Inline.attributed(inlines, size: size))
                    .foregroundStyle(color)
                    .lineSpacing(size * 0.32)
                    .fixedSize(horizontal: false, vertical: true)
            }
        }
    }

    private var paragraphs: [[MDInline]] {
        Readme.parse(text).compactMap { block in
            if case .paragraph(let inlines) = block { inlines } else { nil }
        }
    }
}

/// A chapter's mission patch: its number on a round, stitched patch in the chapter's color once
/// its required lessons are done, and an outline of one until then.
struct PatchView: View {
    let number: Int
    let earned: Bool
    var size: CGFloat = 34

    var body: some View {
        let color = StoryPalette.patch(number)
        ZStack {
            Circle().fill(earned ? color : Color.clear)
            Circle()
                .strokeBorder(earned ? StoryPalette.patchInk.opacity(0.55) : Theme.secondaryText.opacity(0.45),
                              style: StrokeStyle(lineWidth: max(0.8, size * 0.035), dash: [max(1.5, size * 0.07), max(1.5, size * 0.06)]))
                .padding(size * 0.09)
            Circle().strokeBorder(earned ? color : Theme.secondaryText.opacity(0.3), lineWidth: max(1, size * 0.045))
            Image(systemName: "star.fill")
                .font(.system(size: max(4, size * 0.15)))
                .foregroundStyle(earned ? StoryPalette.patchInk.opacity(0.8) : Theme.secondaryText.opacity(0.6))
                .offset(y: -size * 0.27)
            Text("\(number)")
                .font(.system(size: size * 0.36, weight: .bold, design: .rounded))
                .monospacedDigit()
                .foregroundStyle(earned ? StoryPalette.patchInk : Theme.secondaryText)
                .offset(y: size * 0.05)
        }
        .frame(width: size, height: size)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("Chapter \(number) patch, \(earned ? "earned" : "not earned yet")")
    }
}

/// The opening's sky: stars, the curve of the Earth and a satellite crossing above it, a little
/// further along on each page. Still with Reduce Motion.
struct NightSky: View {
    /// The opening's page, which moves the satellite on.
    var page: Int = 0
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        TimelineView(.animation(minimumInterval: 1.0 / 30, paused: reduceMotion)) { timeline in
            let t = reduceMotion ? 0 : timeline.date.timeIntervalSinceReferenceDate
            Canvas { context, size in
                drawStars(in: &context, size: size, time: t)
                drawEarth(in: &context, size: size)
                drawSatellite(in: &context, size: size, time: t)
            }
        }
        .background(StoryPalette.night)
        .accessibilityHidden(true)
    }

    /// The same sky every time: positions from a fixed sequence, not chance.
    private static let stars: [(x: Double, y: Double, r: Double, a: Double, speed: Double)] = {
        var seed: UInt64 = 0x9E37_79B9_7F4A_7C15
        func next() -> Double {
            seed = seed &* 6_364_136_223_846_793_005 &+ 1_442_695_040_888_963_407
            return Double(seed >> 11) / Double(1 << 53)
        }
        return (0..<150).map { _ in (next(), next() * 0.8, 0.5 + next() * 1.1, 0.25 + next() * 0.65, 0.4 + next() * 1.6) }
    }()

    private func drawStars(in context: inout GraphicsContext, size: CGSize, time: Double) {
        for (i, star) in NightSky.stars.enumerated() {
            let twinkle = i % 3 == 0 ? 0.65 + 0.35 * sin(time * star.speed + Double(i)) : 1
            let rect = CGRect(x: star.x * size.width, y: star.y * size.height, width: star.r * 2, height: star.r * 2)
            context.fill(Path(ellipseIn: rect), with: .color(.white.opacity(star.a * twinkle)))
        }
    }

    /// The Earth's center and radius: its curve along the bottom of the window.
    private func earth(_ size: CGSize) -> (center: CGPoint, radius: Double) {
        let radius = max(size.width, size.height) * 1.15
        return (CGPoint(x: size.width * 0.62, y: size.height + radius * 0.9), radius)
    }

    private func drawEarth(in context: inout GraphicsContext, size: CGSize) {
        let (c, r) = earth(size)
        let disc = Path(ellipseIn: CGRect(x: c.x - r, y: c.y - r, width: r * 2, height: r * 2))
        context.fill(disc, with: .color(StoryPalette.earth))
        context.stroke(disc, with: .color(Theme.orbit.opacity(0.14)), lineWidth: 9)
        context.stroke(disc, with: .color(Theme.orbit.opacity(0.55)), lineWidth: 1.2)
    }

    private func drawSatellite(in context: inout GraphicsContext, size: CGSize, time: Double) {
        let (c, r) = earth(size)
        // Low over the Earth, so the orbit stays below the page's text.
        let orbit = r + 48
        // An arc of the orbit over the visible sky, dashed.
        var arc = Path()
        arc.addArc(center: c, radius: orbit, startAngle: .degrees(-140), endAngle: .degrees(-40), clockwise: false)
        context.stroke(arc, with: .color(Theme.orbit.opacity(0.2)), style: StrokeStyle(lineWidth: 1, dash: [3, 7]))
        // A page moves it on by a few degrees; time, slowly, between pages.
        let angle = (-122 + Double(page) * 8 + (time.truncatingRemainder(dividingBy: 240)) / 240 * 8) * .pi / 180
        let p = CGPoint(x: c.x + orbit * cos(angle), y: c.y + orbit * sin(angle))
        var body = context
        body.translateBy(x: p.x, y: p.y)
        body.rotate(by: .radians(angle + .pi / 2))
        body.fill(Path(CGRect(x: -3, y: -3, width: 6, height: 6)), with: .color(.white))
        body.fill(Path(CGRect(x: -15, y: -1.6, width: 10, height: 3.2)), with: .color(Theme.orbit))
        body.fill(Path(CGRect(x: 5, y: -1.6, width: 10, height: 3.2)), with: .color(Theme.orbit))
        context.fill(Path(ellipseIn: CGRect(x: p.x - 14, y: p.y - 14, width: 28, height: 28)), with: .color(Theme.orbit.opacity(0.06)))
    }
}

/// Initials for a badge: the first letters of a name's first two words.
func initials(_ name: String) -> String {
    name.split(whereSeparator: \.isWhitespace).prefix(2).compactMap(\.first).map { String($0).uppercased() }.joined()
}
