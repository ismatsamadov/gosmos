import AppKit
import GosmosKit
import SwiftUI

/// A course's opening, over the whole window: its pages one at a time on the night sky, each
/// with its text, its cards or the learner's badge, and the button that moves on. The last
/// page's button opens the course. It can be skipped, and replayed from Mission Control.
struct OpeningView: View {
    let opening: AppModel.Opening
    @Environment(AppModel.self) private var app
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    /// The badge's name while it is typed; it is saved when the page moves on.
    @ViewState private var name = ""
    @FocusState private var nameFocused: Bool

    var body: some View {
        let pages = app.story(of: opening.course, in: opening.library)?.pages ?? []
        ZStack {
            NightSky(page: opening.page)
                .animation(reduceMotion ? nil : .smooth(duration: 1.2), value: opening.page)
            if pages.indices.contains(opening.page) {
                page(pages[opening.page], count: pages.count)
                    .id(opening.page)
                    // A native text field must leave with its page, rather than stay alive as
                    // an outgoing animated view while the next page is already on screen.
                    .transaction { transaction in
                        transaction.animation = nil
                        transaction.disablesAnimations = true
                    }
            }
        }
        .overlay(alignment: .topTrailing) {
            Button("Skip Intro") { leavePage { app.finishOpening() } }
                .secondaryGlassButton()
                .keyboardShortcut(.cancelAction)
                .help("Skip the opening and start the course (Esc)")
                .padding(.top, 18).padding(.trailing, 22)
        }
        .environment(\.colorScheme, .dark)
        .onAppear { name = app.badgeName }
    }

    private func page(_ page: CourseStory.Page, count: Int) -> some View {
        GeometryReader { geometry in
            let wide = geometry.size.width >= 1000
            let layout = wide
                ? AnyLayout(HStackLayout(alignment: .center, spacing: 56))
                : AnyLayout(VStackLayout(alignment: .leading, spacing: 32))
            ScrollView {
                // One set of controls in either layout. ViewThatFits created two name
                // fields bound to the same focus state, including an off-screen one.
                layout {
                    text(page, count: count)
                        .frame(width: wide ? 470 : nil, alignment: .leading)
                    side(page)
                        .frame(maxWidth: wide ? 520 : .infinity)
                }
                .padding(.horizontal, wide ? 64 : 40)
                .padding(.top, 74)
                .padding(.bottom, 40)
                .frame(maxWidth: 1180)
                .frame(maxWidth: .infinity, minHeight: geometry.size.height)
            }
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    private func text(_ page: CourseStory.Page, count: Int) -> some View {
        VStack(alignment: .leading, spacing: 22) {
            if let eyebrow = page.eyebrow {
                Text(eyebrow.uppercased())
                    .font(.system(size: 11, weight: .medium, design: .monospaced)).tracking(2.2)
                    .foregroundStyle(Theme.orbit)
                    .fixedSize(horizontal: false, vertical: true)
            }
            Text(page.title)
                .font(.system(size: 42, weight: .semibold)).tracking(-1.2)
                .foregroundStyle(.white)
                .fixedSize(horizontal: false, vertical: true)
                .accessibilityAddTraits(.isHeader)
            StoryText(text: page.body, size: 17, color: Theme.onInk, spacing: 14)
            if page.kind == .badge {
                TextField("Your name", text: $name, prompt: Text("What should the team call you?"))
                    .textFieldStyle(.roundedBorder)
                    .controlSize(.large)
                    .font(.system(size: 15))
                    .frame(maxWidth: 340)
                    .focused($nameFocused)
                    // Return is the page's button, the default action, as anywhere on the page.
                    .onSubmit { leavePage { app.advanceOpening() } }
                    .onExitCommand { leavePage { app.finishOpening() } }
                    .onAppear { nameFocused = true }
                    .accessibilityLabel("The name on your badge")
            }
            HStack(spacing: 14) {
                Button { leavePage { app.advanceOpening() } } label: {
                    HStack(spacing: 14) {
                        Text(page.action)
                        Image(systemName: "arrow.right")
                    }
                    .font(.system(size: 14, weight: .semibold)).padding(.horizontal, 8).padding(.vertical, 6)
                }
                .prominentGlassButton()
                .keyboardShortcut(.defaultAction)
                if opening.page > 0 {
                    Button { leavePage { app.openingPageBack() } } label: { Label("Back", systemImage: "chevron.backward") }
                        .buttonStyle(.plain).foregroundStyle(Theme.onInk)
                        .help("The page before")
                }
            }
            .padding(.top, 6)
            PageDots(page: opening.page, count: count).padding(.top, 10)
        }
    }

    @ViewBuilder private func side(_ page: CourseStory.Page) -> some View {
        if page.kind == .badge {
            BadgeView(header: page.eyebrow, name: name.trimmingCharacters(in: .whitespaces), cards: page.cards)
                .rotationEffect(.degrees(reduceMotion ? 0 : -2.5))
                .frame(maxWidth: .infinity)
        } else if !page.cards.isEmpty {
            StoryCardGrid(cards: page.cards)
        }
    }

    /// Commit editing before removing the field. Back and Skip keep the name too.
    private func leavePage(_ action: () -> Void) {
        nameFocused = false
        NSApp.keyWindow?.makeFirstResponder(nil)
        let pages = app.story(of: opening.course, in: opening.library)?.pages ?? []
        if pages.indices.contains(opening.page), pages[opening.page].kind == .badge {
            app.setBadgeName(name)
        }
        action()
    }
}

/// Where the opening is: a dot a page, the current one long.
private struct PageDots: View {
    let page: Int
    let count: Int

    var body: some View {
        HStack(spacing: 6) {
            ForEach(0..<count, id: \.self) { i in
                Capsule().fill(i == page ? Theme.orbit : Theme.onInk.opacity(0.3))
                    .frame(width: i == page ? 18 : 6, height: 6)
            }
        }
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("Page \(page + 1) of \(count)")
    }
}

/// A page's cards: people, places or words, two to a row.
struct StoryCardGrid: View {
    let cards: [CourseStory.Card]

    var body: some View {
        LazyVGrid(columns: [GridItem(.flexible(), spacing: 12, alignment: .top), GridItem(.flexible(), spacing: 12, alignment: .top)],
                  alignment: .leading, spacing: 12) {
            ForEach(cards) { card in
                StoryCardView(card: card)
            }
        }
    }
}

struct StoryCardView: View {
    let card: CourseStory.Card

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack(alignment: .center, spacing: 10) {
                Text(card.label)
                    .font(.system(size: 11, weight: .semibold, design: .monospaced))
                    .foregroundStyle(Theme.orbit)
                    .lineLimit(1)
                    .padding(.horizontal, 8).padding(.vertical, 5)
                    .background(Theme.orbit.opacity(0.12), in: RoundedRectangle(cornerRadius: 7))
                VStack(alignment: .leading, spacing: 2) {
                    Text(card.title).font(.system(size: 13, weight: .semibold)).foregroundStyle(.white)
                        .fixedSize(horizontal: false, vertical: true)
                    if let subtitle = card.subtitle {
                        Text(subtitle).font(.system(size: 11)).foregroundStyle(Theme.onInk)
                            .fixedSize(horizontal: false, vertical: true)
                    }
                }
            }
            if let body = card.body {
                StoryText(text: body, size: 12, color: Theme.onInk.opacity(0.92), spacing: 6)
            }
        }
        .padding(14)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(StoryPalette.panel, in: RoundedRectangle(cornerRadius: 14, style: .continuous))
        .overlay {
            RoundedRectangle(cornerRadius: 14, style: .continuous).strokeBorder(Theme.orbit.opacity(0.16))
        }
        .accessibilityElement(children: .combine)
    }
}
