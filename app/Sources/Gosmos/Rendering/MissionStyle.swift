import SwiftUI

extension View {
    func missionCard(tint: Color? = nil, radius: CGFloat = 14) -> some View {
        modifier(MissionCard(tint: tint, radius: radius))
    }
}

/// A solid surface with a fine outline. The surface is close to the canvas behind it, so with
/// Increase Contrast the outline is what sets the card apart, and it deepens.
private struct MissionCard: ViewModifier {
    let tint: Color?
    let radius: CGFloat
    @Environment(\.colorSchemeContrast) private var contrast

    func body(content: Content) -> some View {
        content
            .background(Theme.surface, in: RoundedRectangle(cornerRadius: radius, style: .continuous))
            .overlay {
                RoundedRectangle(cornerRadius: radius, style: .continuous)
                    .strokeBorder(contrast == .increased ? (tint ?? Theme.outline(contrast)) : (tint?.opacity(0.3) ?? Theme.border), lineWidth: 1)
                    .allowsHitTesting(false)
            }
    }
}

struct MissionIcon: View {
    let symbol: String
    var color: Color = Theme.accent
    var size: CGFloat = 40

    var body: some View {
        Image(systemName: symbol)
            .font(.system(size: size * 0.43, weight: .semibold))
            .symbolRenderingMode(.hierarchical)
            .foregroundStyle(color)
            .frame(width: size, height: size)
            .background(color.opacity(0.10), in: RoundedRectangle(cornerRadius: size * 0.3))
            .accessibilityHidden(true)
    }
}

struct MissionPill: View {
    let text: String
    var symbol: String? = nil
    var color: Color = Theme.accent

    var body: some View {
        HStack(spacing: 5) {
            if let symbol { Image(systemName: symbol) }
            Text(text)
        }
        .font(.system(size: 10, weight: .semibold))
        .foregroundStyle(color)
        .padding(.horizontal, 9)
        .padding(.vertical, 5)
        .background(color.opacity(0.10), in: Capsule())
        .fixedSize()
        .accessibilityElement(children: .combine)
    }
}

struct MissionHeading: View {
    let eyebrow: String
    let title: String
    let subtitle: String
    let symbol: String
    var color: Color = Theme.accent

    var body: some View {
        HStack(alignment: .top, spacing: 16) {
            MissionIcon(symbol: symbol, color: color, size: 48)
            VStack(alignment: .leading, spacing: 7) {
                Text(eyebrow.uppercased())
                    .font(.system(size: 10, weight: .semibold))
                    .tracking(1.5)
                    .foregroundStyle(color)
                Text(title)
                    .font(.system(size: 28, weight: .semibold))
                    .tracking(-0.7)
                    .accessibilityAddTraits(.isHeader)
                Text(subtitle)
                    .font(.system(size: 13))
                    .foregroundStyle(Theme.secondaryText)
                    .fixedSize(horizontal: false, vertical: true)
            }
        }
    }
}

/// A single quiet orbit, drawn with native shapes and still with Reduce Motion.
struct OrbitEmblem: View {
    var color: Color = Theme.accent

    var body: some View {
        ZStack {
            Circle().stroke(color.opacity(0.12), lineWidth: 1).frame(width: 188, height: 188)
            Circle().stroke(color.opacity(0.16), style: StrokeStyle(lineWidth: 1, dash: [2, 6]))
                .frame(width: 144, height: 144)
            Rectangle().fill(color.opacity(0.10)).frame(width: 220, height: 1)
            Rectangle().fill(color.opacity(0.10)).frame(width: 1, height: 210)
            Circle().fill(color.opacity(0.07)).frame(width: 100, height: 100)
            Circle().stroke(color.opacity(0.28), lineWidth: 1).frame(width: 100, height: 100)
            Text("{ }")
                .font(.system(size: 36, weight: .light, design: .monospaced))
                .foregroundStyle(color)
            Ellipse().stroke(color.opacity(0.5), lineWidth: 1)
                .frame(width: 238, height: 88).rotationEffect(.degrees(-32))
            Circle().fill(color).frame(width: 9, height: 9).offset(x: 93, y: -58)
            Circle().fill(Theme.gold).frame(width: 6, height: 6).offset(x: -86, y: 67)
            Image(systemName: "plus").font(.system(size: 10, weight: .light))
                .foregroundStyle(color.opacity(0.6)).offset(x: -88, y: -70)
            Text("GO / EXPLORE")
                .font(.system(size: 8, weight: .medium, design: .monospaced))
                .tracking(2).foregroundStyle(color.opacity(0.7)).offset(y: 115)
        }
        .frame(width: 250, height: 248)
        .accessibilityHidden(true)
    }
}

/// A slim progress track whose empty state is actually empty.
struct MissionProgress: View {
    let fraction: Double
    var color: Color = Theme.accent
    var height: CGFloat = 4
    @Environment(\.colorSchemeContrast) private var contrast

    var body: some View {
        GeometryReader { geometry in
            ZStack(alignment: .leading) {
                Capsule().fill(contrast == .increased ? Theme.secondaryText.opacity(0.35) : color.opacity(0.12))
                Capsule().fill(color).frame(width: geometry.size.width * min(1, max(0, fraction)))
            }
        }
        .frame(height: height)
        .accessibilityElement(children: .ignore)
        .accessibilityValue("\(percent) percent")
    }

    /// Whole percent, rounded down as the labels beside a track count it (`done * 100 / total`):
    /// 29 of 50 is 0.58, and 0.58 * 100 is just under 58 in binary.
    private var percent: Int { Int(min(1, max(0, fraction)) * 100 + 1e-9) }
}

struct GosmosMark: View {
    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 11).fill(Theme.accent.opacity(0.10))
            Circle().stroke(Theme.accent, lineWidth: 1.6).frame(width: 17, height: 17)
            Ellipse().stroke(Theme.accent, lineWidth: 1.4)
                .frame(width: 29, height: 10).rotationEffect(.degrees(-35))
            Circle().fill(Theme.accent).frame(width: 4, height: 4).offset(x: 11, y: -7)
        }
        .frame(width: 36, height: 36)
        .accessibilityHidden(true)
    }
}

/// Gentle pointer feedback for content rows without moving their labels or focus targets.
/// `feedback` false keeps a row still where a click changes nothing, such as a graded answer,
/// without dimming it as disabling would.
struct MissionRowButtonStyle: ButtonStyle {
    var feedback = true

    func makeBody(configuration: Configuration) -> some View {
        MissionRowButton(configuration: configuration, feedback: feedback)
    }
}

private struct MissionRowButton: View {
    let configuration: ButtonStyleConfiguration
    let feedback: Bool
    @ViewState private var hovering = false
    @Environment(\.isEnabled) private var isEnabled

    var body: some View {
        let active = isEnabled && feedback
        configuration.label
            .overlay {
                RoundedRectangle(cornerRadius: 10)
                    .fill(Theme.accent.opacity(active && configuration.isPressed ? 0.10 : active && hovering ? 0.05 : 0))
                    .allowsHitTesting(false)
            }
            .opacity(isEnabled ? 1 : 0.45)
            .onHover { hovering = $0 }
    }
}
