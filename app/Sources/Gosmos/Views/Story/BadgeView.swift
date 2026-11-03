import GosmosKit
import SwiftUI

/// The learner's badge, as the opening's badge page makes it: the page's eyebrow at its head
/// (the organization), the name they gave, the first card's title under it (their role), and
/// the other cards as its lines. A card, in the same colors in either appearance.
struct BadgeView: View {
    let header: String?
    let name: String
    let cards: [CourseStory.Card]
    /// Earned patches, by chapter number, shown along its foot; none in the opening.
    var patches: [Int] = []

    private static let paper = Color(red: 246 / 255, green: 249 / 255, blue: 250 / 255)
    private static let ink = Color(red: 16 / 255, green: 34 / 255, blue: 43 / 255)
    private static let band = Color(red: 11 / 255, green: 60 / 255, blue: 73 / 255)
    private static let label = Color(red: 80 / 255, green: 96 / 255, blue: 106 / 255)

    var body: some View {
        VStack(spacing: 0) {
            ZStack {
                BadgeView.band
                HStack(spacing: 8) {
                    Text((header ?? "").uppercased())
                        .font(.system(size: 13, weight: .bold, design: .rounded)).tracking(3.2)
                        .foregroundStyle(.white)
                        .lineLimit(1).minimumScaleFactor(0.6)
                    Spacer(minLength: 0)
                    OrbitMark().frame(width: 30, height: 22)
                }
                .padding(.horizontal, 18)
                Capsule().fill(StoryPalette.night).frame(width: 46, height: 8).offset(y: -22)
            }
            .frame(height: 64)
            VStack(alignment: .leading, spacing: 14) {
                HStack(alignment: .center, spacing: 14) {
                    ZStack {
                        Circle().fill(BadgeView.band)
                        if name.isEmpty {
                            Image(systemName: "person.fill").font(.system(size: 22)).foregroundStyle(Theme.orbit)
                        } else {
                            Text(initials(name)).font(.system(size: 21, weight: .semibold, design: .rounded)).foregroundStyle(Theme.orbit)
                        }
                    }
                    .frame(width: 58, height: 58)
                    VStack(alignment: .leading, spacing: 4) {
                        Text(name.isEmpty ? "Your name here" : name)
                            .font(.system(size: 20, weight: .semibold)).tracking(-0.4)
                            .foregroundStyle(name.isEmpty ? BadgeView.label : BadgeView.ink)
                            .lineLimit(2).minimumScaleFactor(0.7)
                        if let role = cards.first {
                            Text(role.title).font(.system(size: 12, weight: .medium)).foregroundStyle(BadgeView.band)
                                .fixedSize(horizontal: false, vertical: true)
                        }
                    }
                }
                Rectangle().fill(BadgeView.ink.opacity(0.1)).frame(height: 1)
                VStack(alignment: .leading, spacing: 9) {
                    ForEach(cards.dropFirst()) { card in
                        HStack(alignment: .firstTextBaseline, spacing: 10) {
                            Text(card.label.uppercased())
                                .font(.system(size: 9, weight: .semibold)).tracking(1.2)
                                .foregroundStyle(BadgeView.label)
                                .frame(width: 74, alignment: .leading)
                            Text(card.title).font(.system(size: 12, weight: .medium, design: card.title.first?.isNumber == true || card.title.hasPrefix("No.") ? .monospaced : .default))
                                .foregroundStyle(BadgeView.ink)
                                .fixedSize(horizontal: false, vertical: true)
                        }
                    }
                }
                HStack(spacing: 5) {
                    if patches.isEmpty {
                        Image(systemName: "seal").font(.system(size: 11)).foregroundStyle(BadgeView.label)
                        Text("A patch for every chapter you finish").font(.system(size: 10)).foregroundStyle(BadgeView.label)
                    } else {
                        ForEach(patches.suffix(8), id: \.self) { PatchView(number: $0, earned: true, size: 22) }
                    }
                }
                .padding(.top, 4)
            }
            .padding(18)
        }
        .frame(width: 290)
        .background(BadgeView.paper, in: RoundedRectangle(cornerRadius: 16, style: .continuous))
        .clipShape(RoundedRectangle(cornerRadius: 16, style: .continuous))
        .overlay {
            RoundedRectangle(cornerRadius: 16, style: .continuous).strokeBorder(.white.opacity(0.22), lineWidth: 1)
        }
        .accessibilityElement(children: .ignore)
        .accessibilityLabel(accessibilityText)
    }

    private var accessibilityText: String {
        var parts = ["\(header ?? "") badge", name.isEmpty ? "no name yet" : name]
        parts += cards.prefix(1).map(\.title)
        parts += cards.dropFirst().map { "\($0.label) \($0.title)" }
        return parts.joined(separator: ", ")
    }
}

/// The badge's mark: an orbit, with its highest point, its apogee, marked.
private struct OrbitMark: View {
    var body: some View {
        ZStack {
            Ellipse().stroke(Theme.orbit, lineWidth: 1.4).frame(width: 28, height: 12).rotationEffect(.degrees(-18))
            Circle().fill(Theme.orbit.opacity(0.35)).frame(width: 7, height: 7).offset(x: -6, y: 2)
            Circle().fill(.white).frame(width: 4.5, height: 4.5).offset(x: 13, y: -4.4)
        }
        .accessibilityHidden(true)
    }
}
