import GosmosKit
import SwiftUI

struct DashboardView: View {
    @Environment(AppModel.self) private var app
    private var progress: LearningProgress { app.learningProgress }

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                heading
                hero
                stats
                ViewThatFits(in: .horizontal) {
                    HStack(alignment: .top, spacing: 20) {
                        courses.frame(minWidth: 420, maxWidth: .infinity)
                        supportingCards.frame(width: 280)
                    }
                    VStack(alignment: .leading, spacing: 24) {
                        courses
                        supportingCards
                    }
                }
                Label("At your own pace. All your progress, right here on your Mac.", systemImage: "lock.shield")
                    .font(.system(size: 11)).foregroundStyle(Theme.secondaryText)
                    .frame(maxWidth: .infinity).padding(.top, 2)
            }
            .padding(26).frame(maxWidth: 1180).frame(maxWidth: .infinity)
        }
        .background(Theme.canvas)
        .scrollEdgeEffectStyle(.hard, for: .top)
    }

    private var heading: some View {
        VStack(alignment: .leading, spacing: 7) {
            Text("A little Go. A long way.")
                .font(.system(size: 28, weight: .semibold)).tracking(-0.9)
                .accessibilityAddTraits(.isHeader)
            Text(progress.completedLessons == 0
                 ? "Big ideas start with small steps. Let’s take the first one."
                 : "You’re building something good. Pick up where you left off.")
                .font(.system(size: 13)).foregroundStyle(Theme.secondaryText)
        }
    }

    private var hero: some View {
        // Worked out once: both layouts below show the same lesson.
        let next = app.recommendedLesson.flatMap { ref in app.summary(ref).map { (ref: ref, lesson: $0) } }
        return ViewThatFits(in: .horizontal) {
            HStack(spacing: 16) {
                heroContent(next).frame(minWidth: 290, maxWidth: .infinity, alignment: .leading)
                OrbitEmblem(color: Theme.orbit).scaleEffect(0.64).frame(width: 200, height: 156)
            }
            heroContent(next)
        }
        .padding(24).frame(maxWidth: .infinity, alignment: .leading)
        .background(Theme.ink, in: RoundedRectangle(cornerRadius: 18))
        .overlay {
            RoundedRectangle(cornerRadius: 18).strokeBorder(Theme.orbit.opacity(0.15))
                .allowsHitTesting(false)
        }
        // This panel stays dark in either appearance.
        .environment(\.colorScheme, .dark)
    }

    private func heroContent(_ next: (ref: LessonRef, lesson: LessonSummary)?) -> some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(spacing: 7) {
                Circle().fill(Theme.orbit).frame(width: 5, height: 5)
                Text(progress.completedLessons == 0 ? "YOUR FIRST SMALL STEP" : "UP NEXT ON YOUR JOURNEY")
                    .font(.system(size: 9, weight: .semibold)).tracking(1.6).foregroundStyle(Theme.orbit)
            }
            if case let (ref, lesson)? = next {
                VStack(alignment: .leading, spacing: 9) {
                    Text(lesson.title).font(.system(size: 30, weight: .semibold)).tracking(-0.7)
                        .foregroundStyle(.white).fixedSize(horizontal: false, vertical: true)
                        .accessibilityAddTraits(.isHeader)
                    if let catalog = app.library(ref.library)?.catalog {
                        Text([catalog.course(of: ref.lesson)?.title, catalog.chapter(of: ref.lesson)?.title]
                            .compactMap { $0 }.joined(separator: "  /  "))
                            .font(.system(size: 12)).foregroundStyle(Theme.onInk)
                            .fixedSize(horizontal: false, vertical: true)
                    }
                }
                ViewThatFits(in: .horizontal) {
                    HStack(spacing: 16) {
                        continueButton(ref)
                        lessonFacts(lesson)
                    }
                    VStack(alignment: .leading, spacing: 12) {
                        lessonFacts(lesson)
                        continueButton(ref)
                    }
                }
                .padding(.top, 3)
            } else {
                // No lesson to recommend: every lesson is done, or no library on the shelf has a
                // lesson to read (none opened, or none that can be read).
                let noLessons = app.shelf.allSatisfy { $0.catalog.lessonCount == 0 }
                Text(noLessons ? "Make room for a new idea." : "Look how far you’ve come.")
                    .font(.system(size: 30, weight: .semibold)).tracking(-0.7).foregroundStyle(.white)
                Text(noLessons ? "Open a library and find your first lesson." : "Keep your knowledge fresh with a little practice.")
                    .font(.system(size: 13)).foregroundStyle(Theme.onInk)
                if noLessons {
                    Button("Open a Library…", systemImage: "plus") { app.chooseLibrary() }
                        .prominentGlassButton().disabled(app.isAddingLibrary)
                } else {
                    Button("Practice with Review", systemImage: "arrow.trianglehead.2.clockwise.rotate.90") { app.showReview() }
                        .prominentGlassButton()
                }
            }
        }
    }

    private func continueButton(_ ref: LessonRef) -> some View {
        Button { app.open(ref) } label: {
            HStack(spacing: 16) {
                Text(progress.completedLessons == 0 ? "Start learning" : "Continue learning")
                Image(systemName: "arrow.right")
            }
            .font(.system(size: 12, weight: .semibold)).padding(.horizontal, 6).padding(.vertical, 5)
        }
        .prominentGlassButton()
    }

    private func lessonFacts(_ lesson: LessonSummary) -> some View {
        HStack(spacing: 10) {
            Label(lesson.type.label, systemImage: lesson.type.symbol)
            Text("·")
            Label("\(LearningProgress.lessonXP) XP", systemImage: "bolt.fill")
        }
        .font(.system(size: 11)).foregroundStyle(Theme.onInk).fixedSize()
    }

    private var stats: some View {
        // Each cell is made once, then laid out in a row, or in a grid when the row doesn't fit.
        let progress = self.progress
        let xp = stat("Total experience", value: progress.xp.formatted(), suffix: "XP", symbol: "bolt", color: Theme.gold)
        let lessons = stat("Lessons complete", value: "\(progress.completedLessons)", suffix: "lessons", symbol: "checkmark.circle", color: Theme.accent)
        let rank = stat("Current rank", value: progress.rank, suffix: "Level \(progress.level)", symbol: "sparkles", color: Theme.violet)
        return ViewThatFits(in: .horizontal) {
            HStack(spacing: 0) {
                xp
                statDivider
                lessons
                statDivider
                rank
                statDivider
                reviewStat
            }
            Grid(alignment: .leading, horizontalSpacing: 0, verticalSpacing: 20) {
                GridRow {
                    xp
                    lessons
                }
                GridRow {
                    rank
                    reviewStat
                }
            }
        }
        .padding(.vertical, 19).missionCard()
    }

    private var statDivider: some View {
        Rectangle().fill(Theme.border).frame(width: 1, height: 37).accessibilityHidden(true)
    }

    private var reviewStat: some View {
        Button { app.showReview() } label: {
            stat("Ready to review", value: app.dueCount == 0 ? "All clear" : "\(app.dueCount)",
                 suffix: app.dueCount == 0 ? "" : "due", symbol: "arrow.trianglehead.2.clockwise.rotate.90", color: Theme.success)
                .contentShape(Rectangle())
        }
        .buttonStyle(.plain).help("Open Review (⇧⌘Y)")
    }

    private func stat(_ title: String, value: String, suffix: String, symbol: String, color: Color) -> some View {
        VStack(alignment: .leading, spacing: 10) {
            Label(title, systemImage: symbol).font(.system(size: 11))
                .foregroundStyle(Theme.secondaryText).labelStyle(.titleAndIcon)
            HStack(alignment: .firstTextBaseline, spacing: 5) {
                Text(value).font(.system(size: 22, weight: .semibold)).tracking(-0.5)
                    .foregroundStyle(.primary).monospacedDigit()
                if !suffix.isEmpty {
                    Text(suffix).font(.system(size: 10, weight: .medium)).foregroundStyle(color)
                }
            }
            .lineLimit(1).minimumScaleFactor(0.8)
        }
        .padding(.horizontal, 20).frame(minWidth: 132, maxWidth: .infinity, alignment: .leading)
        .accessibilityElement(children: .combine)
    }

    private var courses: some View {
        VStack(alignment: .leading, spacing: 14) {
            HStack(alignment: .firstTextBaseline) {
                Text("Your learning path").font(.system(size: 16, weight: .semibold)).accessibilityAddTraits(.isHeader)
                Spacer()
                Button { app.chooseLibrary() } label: { Label("Add library", systemImage: "plus") }
                    .font(.system(size: 11, weight: .medium)).buttonStyle(.plain).foregroundStyle(Theme.accent)
                    .disabled(app.isAddingLibrary).frame(minHeight: 24)
            }
            if app.shelf.allSatisfy({ $0.catalog.courses.isEmpty }) {
                ContentUnavailableView {
                    Label("Your next adventure awaits", systemImage: "books.vertical")
                } description: {
                    Text("Open a course library to begin learning.")
                } actions: {
                    Button("Open Library…") { app.chooseLibrary() }.secondaryGlassButton().disabled(app.isAddingLibrary)
                }
                .frame(maxWidth: .infinity).missionCard()
            } else {
                let featured = app.shelf.reduce(0) { $0 + $1.catalog.courses.count } == 1
                LazyVGrid(columns: featured ? [GridItem(.flexible())] : [GridItem(.adaptive(minimum: 290), spacing: 14)], spacing: 14) {
                    ForEach(app.shelf) { library in
                        ForEach(library.catalog.courses) { course in
                            CourseMissionCard(course: course, library: library.id, featured: featured)
                        }
                    }
                }
            }
        }
    }

    private var supportingCards: some View {
        VStack(alignment: .leading, spacing: 14) {
            HStack {
                Text("Milestones").font(.system(size: 16, weight: .semibold)).accessibilityAddTraits(.isHeader)
                Spacer()
                Text("\(progress.achievements.filter(\.isUnlocked).count) / \(progress.achievements.count)")
                    .font(.system(size: 11, weight: .medium).monospacedDigit()).foregroundStyle(Theme.secondaryText)
            }
            .frame(minHeight: 24)
            VStack(alignment: .leading, spacing: 18) {
                Text("Every small win counts.").font(.system(size: 12)).foregroundStyle(Theme.secondaryText)
                LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible())], spacing: 18) {
                    ForEach(progress.achievements) { achievement in AchievementCard(achievement: achievement) }
                }
            }
            .padding(18).missionCard()
            if app.docs != nil {
                Button { app.openDocumentation() } label: {
                    HStack(spacing: 12) {
                        MissionIcon(symbol: "text.book.closed", color: Theme.violet, size: 34)
                        VStack(alignment: .leading, spacing: 4) {
                            Text("A little help along the way").font(.system(size: 11, weight: .semibold))
                            Text("Explore the Go field guide").font(.system(size: 10)).foregroundStyle(Theme.secondaryText)
                        }
                        Spacer(minLength: 0)
                        Image(systemName: "arrow.up.right").font(.system(size: 10)).foregroundStyle(Theme.violet)
                    }
                    .padding(16).frame(maxWidth: .infinity, alignment: .leading)
                    .missionCard().contentShape(RoundedRectangle(cornerRadius: 14))
                }
                .buttonStyle(MissionRowButtonStyle())
            }
        }
    }
}

private struct CourseMissionCard: View {
    let course: Course
    let library: String
    var featured = false
    @Environment(AppModel.self) private var app

    /// Keep the path near the next lesson, with context from the preceding chapter; the last
    /// chapters once every lesson is done.
    private func visibleChapters(around upNext: LessonSummary?) -> [Chapter] {
        let chapters = course.chapters
        let current = upNext.flatMap { next in chapters.firstIndex { $0.lessons.contains { $0.id == next.id } } }
            ?? max(0, chapters.count - 1)
        let start = min(max(0, current - 1), max(0, chapters.count - 3))
        return Array(chapters.dropFirst(start).prefix(3))
    }

    var body: some View {
        let progress = app.progress(of: course, in: library)
        let complete = progress.total > 0 && progress.done == progress.total
        // Worked out once for the card and its rows: the next lesson not done (nil once every
        // lesson is), and what the course action opens.
        let upNext = app.firstUnfinished(course.lessons, in: library)
        let next = upNext ?? course.lessons.first
        let chapters = featured ? visibleChapters(around: upNext) : []
        VStack(alignment: .leading, spacing: 0) {
            VStack(alignment: .leading, spacing: 20) {
                HStack(alignment: .top, spacing: 12) {
                    MissionIcon(symbol: complete ? "checkmark.seal" : course.isGuidedProject ? "hammer" : "curlybraces", size: 40)
                    VStack(alignment: .leading, spacing: 5) {
                        Text(course.title).font(.system(size: 18, weight: .semibold)).tracking(-0.4)
                            .fixedSize(horizontal: false, vertical: true)
                        Text("\(course.chapters.count) chapters  ·  \(course.lessons.count) lessons")
                            .font(.system(size: 11)).foregroundStyle(Theme.secondaryText)
                    }
                    Spacer(minLength: 0)
                    MissionPill(text: complete ? "Complete" : course.isGuidedProject ? "Project" : "Course", color: complete ? Theme.success : Theme.secondaryText)
                }
                VStack(spacing: 9) {
                    HStack {
                        Text("\(progress.done) of \(progress.total) required lessons").font(.system(size: 11)).foregroundStyle(Theme.secondaryText)
                        Spacer()
                        Text(progress.total == 0 ? "0%" : "\(progress.done * 100 / progress.total)%")
                            .font(.system(size: 11, weight: .semibold).monospacedDigit()).foregroundStyle(Theme.accent)
                    }
                    MissionProgress(fraction: Double(progress.done) / Double(max(1, progress.total)))
                        .accessibilityLabel("\(course.title) progress")
                }
            }
            .padding(20)
            if !chapters.isEmpty {
                Rectangle().fill(Theme.border).frame(height: 1)
                VStack(spacing: 0) {
                    ForEach(chapters) { chapter in
                        chapterRow(chapter, upNext: upNext)
                        if chapter.id != chapters.last?.id {
                            Rectangle().fill(Theme.border.opacity(0.6)).frame(height: 1).padding(.leading, 42)
                        }
                    }
                }
                .padding(.horizontal, 20).padding(.vertical, 4)
            }
            Button {
                if let next { app.open(LessonRef(library: library, lesson: next.id)) }
            } label: {
                HStack {
                    // "Revisit" only when nothing is left: with its required lessons done, a
                    // course can still have an optional lesson to open next.
                    Text(upNext == nil ? "Revisit course" : progress.done == 0 ? "Begin course" : "Continue course")
                    Spacer()
                    Image(systemName: "arrow.right")
                }
                .font(.system(size: 12, weight: .semibold)).foregroundStyle(Theme.accent)
                .padding(.horizontal, 20).padding(.vertical, 15)
                .background(Theme.accent.opacity(0.045)).contentShape(Rectangle())
            }
            .buttonStyle(MissionRowButtonStyle()).disabled(next == nil)
        }
        .frame(maxWidth: .infinity, alignment: .leading).missionCard()
        .clipShape(RoundedRectangle(cornerRadius: 14))
    }

    private func chapterRow(_ chapter: Chapter, upNext: LessonSummary?) -> some View {
        let done = app.doneCount(chapter, in: library)
        let complete = !chapter.lessons.isEmpty && done == chapter.lessons.count
        // The chapter of the next lesson; none once the whole course is done.
        let current = upNext.map { next in chapter.lessons.contains { $0.id == next.id } } ?? false
        let index = (course.chapters.firstIndex { $0.id == chapter.id } ?? 0) + 1
        let target = app.firstUnfinished(chapter.lessons, in: library) ?? chapter.lessons.first
        return Button {
            if let target { app.open(LessonRef(library: library, lesson: target.id)) }
        } label: {
            HStack(spacing: 12) {
                ZStack {
                    RoundedRectangle(cornerRadius: 9).fill(current ? Theme.accent.opacity(0.10) : Theme.inset)
                    if complete {
                        Image(systemName: "checkmark").font(.system(size: 11, weight: .semibold)).foregroundStyle(Theme.success)
                    } else {
                        Text(String(format: "%02d", index)).font(.system(size: 10, weight: .medium, design: .monospaced))
                            .foregroundStyle(current ? Theme.accent : Theme.secondaryText)
                    }
                }
                .frame(width: 30, height: 30)
                VStack(alignment: .leading, spacing: 4) {
                    Text(chapter.title).font(.system(size: 12, weight: current ? .semibold : .medium)).foregroundStyle(.primary).lineLimit(2)
                    Text("\(done) of \(chapter.lessons.count) lessons complete").font(.system(size: 10)).foregroundStyle(Theme.secondaryText)
                }
                Spacer(minLength: 4)
                if current { MissionPill(text: "Next up") }
                Image(systemName: "chevron.right").font(.system(size: 9, weight: .medium)).foregroundStyle(Theme.secondaryText)
            }
            .padding(.vertical, 13).contentShape(Rectangle())
        }
        .buttonStyle(MissionRowButtonStyle()).disabled(target == nil)
        // What the row shows, the Next up pill included; the button opens a lesson of the chapter.
        .accessibilityLabel("Chapter \(index), \(chapter.title)\(current ? ", next up" : ""), \(done) of \(chapter.lessons.count) lessons complete")
        .accessibilityHint(target.map { "Opens \($0.title)" } ?? "")
    }
}

private struct AchievementCard: View {
    let achievement: LearningProgress.Achievement

    var body: some View {
        VStack(spacing: 7) {
            ZStack(alignment: .bottomTrailing) {
                MissionIcon(symbol: achievement.symbol, color: achievement.isUnlocked ? Theme.gold : Theme.secondaryText, size: 38)
                if achievement.isUnlocked {
                    Image(systemName: "checkmark.circle.fill").font(.system(size: 11)).foregroundStyle(Theme.gold)
                        .background(Theme.surface, in: Circle()).offset(x: 3, y: 2)
                }
            }
            Text(achievement.title).font(.system(size: 11, weight: .medium))
                .foregroundStyle(achievement.isUnlocked ? .primary : Theme.secondaryText).lineLimit(1)
            Text(achievement.isUnlocked ? "Unlocked" : "\(min(achievement.value, achievement.target)) / \(achievement.target)")
                .font(.system(size: 9, design: .monospaced)).foregroundStyle(achievement.isUnlocked ? Theme.gold : Theme.secondaryText)
        }
        .frame(maxWidth: .infinity).help(achievement.detail)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("\(achievement.title). \(achievement.detail). \(achievement.isUnlocked ? "Unlocked" : "\(achievement.value) of \(achievement.target)")")
    }
}

struct SidebarBrand: View {
    var body: some View {
        HStack(spacing: 10) {
            GosmosMark()
            VStack(alignment: .leading, spacing: 3) {
                Text("gosmos").font(.system(size: 23, weight: .semibold, design: .rounded)).tracking(-0.8)
                Text("A little Go. A long way.").font(.system(size: 10)).foregroundStyle(Theme.secondaryText)
            }
            Spacer(minLength: 0)
        }
        .padding(.horizontal, 20).padding(.top, 18).padding(.bottom, 24).background(Theme.sidebar)
        .accessibilityElement(children: .combine)
    }
}

struct SidebarProgress: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        let progress = app.learningProgress
        Button { app.showDashboard() } label: {
            VStack(alignment: .leading, spacing: 12) {
                HStack(spacing: 9) {
                    MissionIcon(symbol: "sparkles", color: Theme.violet, size: 30)
                    VStack(alignment: .leading, spacing: 3) {
                        Text("\(progress.rank) · Level \(progress.level)").font(.system(size: 11, weight: .semibold))
                        Text("\(progress.xp.formatted()) XP earned").font(.system(size: 10)).foregroundStyle(Theme.secondaryText)
                    }
                    Spacer(minLength: 0)
                    Image(systemName: "chevron.right").font(.system(size: 9)).foregroundStyle(Theme.secondaryText)
                }
                VStack(alignment: .leading, spacing: 7) {
                    MissionProgress(fraction: progress.levelFraction, color: Theme.violet)
                        .accessibilityLabel("Progress to level \(progress.level + 1)")
                    HStack {
                        Text("Next: Level \(progress.level + 1)")
                        Spacer(minLength: 0)
                        Text("\(progress.xpToNextLevel) XP to go")
                    }
                    .font(.system(size: 9)).foregroundStyle(Theme.secondaryText)
                }
            }
            .padding(14).missionCard(radius: 12).contentShape(RoundedRectangle(cornerRadius: 12))
        }
        .buttonStyle(.plain).padding(12).background(Theme.sidebar)
        .accessibilityLabel("Level \(progress.level), \(progress.rank), \(progress.xp) XP. Open Mission Control.")
    }
}
