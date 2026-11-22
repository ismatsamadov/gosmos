import GosmosKit
import Observation

/// A place Back and Forward return to: Mission Control, the field guide, a lesson, or a place
/// in the Go documentation.
enum Place: Hashable {
    case dashboard
    case fieldGuide
    case lesson(LessonRef)
    case doc(DocLink)
}

/// Places visited, for the toolbar's back and forward buttons.
struct VisitHistory {
    private(set) var back: [Place] = []
    private(set) var forward: [Place] = []

    mutating func visit(from current: Place?) {
        if let current, back.last != current { back.append(current) }
        if back.count > 200 { back.removeFirst() }
        forward.removeAll()
    }

    /// An entry for the place on screen is no step: it is left when you go to a place the
    /// history doesn't keep (the review, a redo) and come straight back.
    mutating func goBack(from current: Place?) -> Place? {
        while let last = back.last, last == current { back.removeLast() }
        guard let target = back.popLast() else { return nil }
        if let current { forward.append(current) }
        return target
    }

    mutating func goForward(from current: Place?) -> Place? {
        while let next = forward.last, next == current { forward.removeLast() }
        guard let target = forward.popLast() else { return nil }
        if let current { back.append(current) }
        return target
    }

    func canGoBack(from current: Place?) -> Bool { back.contains { $0 != current } }
    func canGoForward(from current: Place?) -> Bool { forward.contains { $0 != current } }
}

/// Where you are, for Back and Forward: the place on screen, or about to be, and the places
/// visited before and after it. AppModel goes to places and leaves them; this keeps what the
/// history knows of them.
@Observable @MainActor
final class NavigationModel {
    /// A place on screen as the history knows it. The documentation's is the place it shows at
    /// the time.
    enum Screen: Equatable {
        case dashboard
        case fieldGuide
        case documentation
        case lesson(LessonRef)
    }

    private(set) var history = VisitHistory()
    /// The place on screen, or about to be while a lesson loads or leaves; nil for the review
    /// and a redo, which the history doesn't keep.
    @ObservationIgnored var screen: Screen?
    /// Where the documentation is: its page, and the place on it. AppModel gives it.
    @ObservationIgnored var docLocation: @MainActor () -> DocLink? = { nil }

    /// The lesson that is the place on screen, or about to be while it loads; nil for Mission
    /// Control, the field guide, the documentation, the review and a redo.
    var lessonTarget: LessonRef? {
        if case .lesson(let ref)? = screen { ref } else { nil }
    }

    /// The place on screen as the back and forward history knows it: Mission Control, the
    /// documentation's page, or the lesson, also while it still loads, so Back or Forward
    /// meanwhile keeps it; the review and a redo are not in it.
    var visited: Place? {
        switch screen {
        case .dashboard?: .dashboard
        case .fieldGuide?: .fieldGuide
        case .documentation?: docLocation().map(Place.doc)
        case .lesson(let ref)?: .lesson(ref)
        case nil: nil
        }
    }

    var canGoBack: Bool { history.canGoBack(from: visited) }
    var canGoForward: Bool { history.canGoForward(from: visited) }

    /// Puts the place on screen in the back history as you go on from it; the forward history
    /// goes.
    func visit() { history.visit(from: visited) }

    /// The place Back returns to, with the one on screen put in the forward history; nil when
    /// there is none.
    func goBack() -> Place? { history.goBack(from: visited) }

    /// The place Forward returns to, with the one on screen put in the back history; nil when
    /// there is none.
    func goForward() -> Place? { history.goForward(from: visited) }
}

// MARK: - Back and Forward

extension AppModel {
    var canGoBack: Bool { navigation.canGoBack }
    var canGoForward: Bool { navigation.canGoForward }

    func goBack() {
        if let place = navigation.goBack() { go(to: place) }
    }

    func goForward() {
        if let place = navigation.goForward() { go(to: place) }
    }

    private func go(to place: Place) {
        switch place {
        case .dashboard: showDashboard(remember: false)
        case .fieldGuide: showFieldGuide(remember: false)
        case .lesson(let ref): open(ref, remember: false)
        case .doc(let link): openDocumentation(link, remember: false)
        }
    }
}
