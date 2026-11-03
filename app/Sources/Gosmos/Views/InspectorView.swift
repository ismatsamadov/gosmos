import GosmosKit
import SwiftUI

/// The inspector: your note for the lesson, and its facts.
struct InspectorView: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        if let lesson = app.lesson {
            InspectorContent(lesson: lesson)
        } else {
            VStack(spacing: 16) {
                MissionIcon(symbol: "sparkles", color: Theme.violet, size: 52)
                Text("Your journey")
                    .font(.system(size: 21, weight: .semibold))
                Text("Level \(app.learningProgress.level) · \(app.learningProgress.rank)")
                    .font(.system(size: 13, weight: .semibold)).foregroundStyle(Theme.violet)
                Text("Open a lesson to keep notes, explore its topics, and plan your next small win.")
                    .font(.system(size: 13)).foregroundStyle(Theme.secondaryText)
                    .multilineTextAlignment(.center)
            }
            .padding(24)
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Theme.canvas)
        }
    }
}

struct InspectorContent: View {
    let lesson: LessonModel
    @Environment(AppModel.self) private var app

    var body: some View {
        let detail = lesson.detail
        Form {
            Section {
                HStack(spacing: 12) {
                    MissionIcon(symbol: lesson.isDone ? "checkmark.seal.fill" : "bolt.fill", color: lesson.isDone ? Theme.success : Theme.gold, size: 34)
                    VStack(alignment: .leading, spacing: 4) {
                        Text(lesson.isReview ? "Practice mission" : lesson.isDone ? "Mission complete" : "Your next small win")
                            .font(.system(size: 12, weight: .semibold))
                        Text(lesson.isReview ? "Strengthen what you know" : "100 XP · \(detail.type.label)")
                            .font(.system(size: 11)).foregroundStyle(Theme.secondaryText)
                    }
                }
                .padding(.vertical, 6)
            }
            Section("Note") {
                if lesson.isReview {
                    // A redo is done without help, the note included; nothing typed here could be kept.
                    Text("A redo is done without your note. It is there when you open the lesson.")
                        .font(.caption)
                        .foregroundStyle(Theme.secondaryText)
                        .fixedSize(horizontal: false, vertical: true)
                } else {
                    NoteEditor(lesson: lesson, placeholder: "Write a note for this lesson…", fontSize: 13,
                               accessibilityLabel: "Note for this lesson")
                        .frame(minHeight: 160)
                    Text("Saved as you type.")
                        .font(.caption)
                        .foregroundStyle(Theme.secondaryText)
                }
            }
            Section("Lesson") {
                LabeledContent("Type") { Label(detail.type.label, systemImage: detail.type.symbol).lineLimit(1) }
                LabeledContent("Difficulty") { DifficultyView(level: detail.difficulty) }
                LabeledContent("Course") { FactText(detail.courseTitle) }
                LabeledContent("Chapter") { FactText(detail.chapterTitle) }
                LabeledContent("Status") { FactText(app.isDone(detail.ref) ? "Done" : detail.summary.required ? "Not done" : "Optional") }
            }
            if !detail.topics.isEmpty {
                Section("Topics") {
                    TopicTags(topics: detail.topics)
                }
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .background(Theme.canvas)
    }
}

/// A fact's value in the inspector: one line, cut with an ellipsis when the column is narrow,
/// and whole in its help tag.
struct FactText: View {
    let text: String
    init(_ text: String) { self.text = text }

    var body: some View {
        Text(text)
            .foregroundStyle(Theme.secondaryText)
            .lineLimit(1)
            .truncationMode(.tail)
            .help(text)
    }
}

/// Difficulty as five dots; lessons rate it from 1 (easiest) to 10.
struct DifficultyView: View {
    let level: Int

    /// 1 and 2 are one dot, 9 and 10 five.
    var dots: Int { min(5, max(1, (level + 1) / 2)) }

    var body: some View {
        HStack(spacing: 3) {
            ForEach(0..<5, id: \.self) { i in
                Circle()
                    .fill(i < dots ? AnyShapeStyle(.tint) : AnyShapeStyle(.quaternary))
                    .frame(width: 7, height: 7)
            }
        }
        .help("Difficulty \(level) of 10")
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("Difficulty \(level) of 10")
    }
}

struct TopicTags: View {
    let topics: [String]

    var body: some View {
        FlowLayout(spacing: 6) {
            ForEach(topics, id: \.self) { topic in
                Text(topic)
                    .font(.system(size: 12))
                    .padding(.horizontal, 9)
                    .padding(.vertical, 4)
                    .background(.tint.opacity(0.12), in: Capsule())
                    .foregroundStyle(.tint)
            }
        }
    }
}

/// Lays children out in rows, wrapping as needed.
struct FlowLayout: Layout {
    var spacing: CGFloat = 6

    func sizeThatFits(proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) -> CGSize {
        let width = proposal.width ?? 260
        var x: CGFloat = 0, y: CGFloat = 0, row: CGFloat = 0, widest: CGFloat = 0
        for s in subviews {
            let size = s.sizeThatFits(.unspecified)
            if x > 0 && x + size.width > width {
                y += row + spacing
                x = 0
                row = 0
            }
            x += size.width + spacing
            row = max(row, size.height)
            widest = max(widest, x - spacing)
        }
        return CGSize(width: min(widest, width), height: y + row)
    }

    func placeSubviews(in bounds: CGRect, proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) {
        var x = bounds.minX, y = bounds.minY, row: CGFloat = 0
        for s in subviews {
            let size = s.sizeThatFits(.unspecified)
            if x > bounds.minX && x + size.width > bounds.maxX {
                y += row + spacing
                x = bounds.minX
                row = 0
            }
            s.place(at: CGPoint(x: x, y: y), proposal: ProposedViewSize(size))
            x += size.width + spacing
            row = max(row, size.height)
        }
    }
}

/// The lesson's note. It reads the note on its own, so typing redraws the editor and not the
/// lesson around it.
struct NoteEditor: View {
    @Bindable var lesson: LessonModel
    let placeholder: String
    let fontSize: CGFloat
    let accessibilityLabel: String

    var body: some View {
        TextEditor(text: $lesson.note)
            .font(.system(size: fontSize))
            .scrollContentBackground(.hidden)
            .overlay(alignment: .topLeading) {
                if lesson.note.isEmpty {
                    Text(placeholder)
                        .font(.system(size: fontSize))
                        .foregroundStyle(Theme.secondaryText)
                        .padding(.leading, 5)
                        .allowsHitTesting(false)
                }
            }
            .accessibilityLabel(accessibilityLabel)
    }
}
