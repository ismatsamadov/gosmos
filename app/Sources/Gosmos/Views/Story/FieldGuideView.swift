import GosmosKit
import SwiftUI

/// The field guide: every word of the courses' worlds with its meaning in plain words, a
/// section a course, and a field that narrows the list as you type.
struct FieldGuideView: View {
    @Environment(AppModel.self) private var app
    @ViewState private var filter = ""
    @FocusState private var filterFocused: Bool

    private struct Section: Identifiable {
        let id: String
        let title: String
        let terms: [CourseStory.Term]
    }

    var body: some View {
        let sections = self.sections
        ScrollViewReader { proxy in
            ScrollView {
                VStack(alignment: .leading, spacing: 22) {
                    MissionHeading(eyebrow: "Field Guide", title: "The words of the job, in plain words",
                                   subtitle: "Every word a lesson uses about space, radio or the ops room. In a lesson, the first time a word comes up it links here.",
                                   symbol: "binoculars", color: Theme.violet)
                    TextField("Find a word", text: $filter, prompt: Text("Find a word"))
                        .textFieldStyle(.roundedBorder)
                        .frame(maxWidth: 320)
                        .focused($filterFocused)
                        .onExitCommand { filter = "" }
                        .accessibilityLabel("Find a word in the field guide")
                    if sections.allSatisfy({ $0.terms.isEmpty }) {
                        ContentUnavailableView.search(text: filter).frame(maxWidth: .infinity)
                    }
                    ForEach(sections.filter { !$0.terms.isEmpty }) { section in
                        VStack(alignment: .leading, spacing: 0) {
                            if sections.count > 1 {
                                Text(section.title).font(.system(size: 13, weight: .semibold)).foregroundStyle(Theme.secondaryText)
                                    .padding(.bottom, 10).accessibilityAddTraits(.isHeader)
                            }
                            VStack(alignment: .leading, spacing: 0) {
                                ForEach(Array(section.terms.enumerated()), id: \.element.id) { i, term in
                                    TermRow(term: term, highlighted: term.term == app.fieldGuideFocus)
                                        .id(section.id + "/" + term.term)
                                    if i < section.terms.count - 1 {
                                        Rectangle().fill(Theme.border.opacity(0.7)).frame(height: 1).padding(.leading, 18)
                                    }
                                }
                            }
                            .missionCard()
                        }
                    }
                }
                .padding(26).frame(maxWidth: 820, alignment: .leading).frame(maxWidth: .infinity)
            }
            .background(Theme.canvas)
            .scrollEdgeEffectStyle(.hard, for: .top)
            .task(id: app.fieldGuideFocus) {
                guard let focus = app.fieldGuideFocus, let section = sections.first(where: { $0.terms.contains { $0.term == focus } }) else { return }
                try? await Task.sleep(for: .milliseconds(80))
                proxy.scrollTo(section.id + "/" + focus, anchor: .center)
            }
        }
    }

    /// Each course's terms that match the filter, in the shelf's order.
    private var sections: [Section] {
        let query = filter.trimmingCharacters(in: .whitespaces)
        var out: [Section] = []
        for library in app.shelf {
            for course in library.catalog.courses {
                guard let terms = library.stories[course.id]?.terms, !terms.isEmpty else { continue }
                let shown = query.isEmpty ? terms : terms.filter {
                    $0.term.localizedCaseInsensitiveContains(query) || $0.meaning.localizedCaseInsensitiveContains(query)
                }
                out.append(Section(id: library.id + "/" + course.id, title: course.title, terms: shown))
            }
        }
        return out
    }
}

private struct TermRow: View {
    let term: CourseStory.Term
    var highlighted = false

    var body: some View {
        HStack(alignment: .firstTextBaseline, spacing: 18) {
            Text(term.term)
                .font(.system(size: 14, weight: .semibold))
                .frame(width: 150, alignment: .leading)
                .fixedSize(horizontal: false, vertical: true)
            StoryText(text: term.meaning, size: 13, color: .primary, spacing: 6)
                .frame(maxWidth: .infinity, alignment: .leading)
        }
        .padding(.horizontal, 18).padding(.vertical, 13)
        .background(highlighted ? Theme.violet.opacity(0.08) : .clear)
        .textSelection(.enabled)
        .accessibilityElement(children: .combine)
    }
}

/// A term's meaning over the lesson, after a click on its link: the word, what it means, and
/// the way to the whole field guide.
struct TermCard: View {
    let term: CourseStory.Term
    @Environment(AppModel.self) private var app

    var body: some View {
        HStack(alignment: .top, spacing: 14) {
            MissionIcon(symbol: "binoculars", color: Theme.violet, size: 34)
            VStack(alignment: .leading, spacing: 6) {
                Text(term.term).font(.system(size: 15, weight: .semibold))
                StoryText(text: term.meaning, size: 13, color: .primary, spacing: 6)
                Button("Open the Field Guide") { app.showFieldGuide(focus: term.term) }
                    .buttonStyle(.plain).font(.system(size: 12, weight: .medium)).foregroundStyle(Theme.violet)
                    .padding(.top, 2)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            Button { app.shownTerm = nil } label: { Image(systemName: "xmark").font(.system(size: 11, weight: .semibold)) }
                .buttonStyle(.plain).foregroundStyle(Theme.secondaryText)
                .keyboardShortcut(.cancelAction)
                .help("Close (Esc)")
                .accessibilityLabel("Close")
        }
        .padding(16)
        .frame(maxWidth: Theme.readableWidth)
        .missionCard(tint: Theme.violet)
        .padding(.horizontal, 24).padding(.bottom, 18)
        .textSelection(.enabled)
        .accessibilityElement(children: .contain)
    }
}
