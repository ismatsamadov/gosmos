import GosmosKit
import Testing
@testable import Gosmos

/// Back and Forward: the history's rules, and the place on screen the navigation model keeps
/// them from.
@Suite @MainActor struct NavigationModelTests {
    private let first = Place.lesson(LessonRef(library: "library", lesson: "first"))
    private let second = Place.lesson(LessonRef(library: "library", lesson: "second"))

    @Test func aVisitKeepsThePlaceLeftOnceAndDropsTheForwardHistory() {
        var history = VisitHistory()
        history.visit(from: .dashboard)
        history.visit(from: first)
        history.visit(from: first)
        #expect(history.back == [.dashboard, first])
        #expect(history.goBack(from: second) == first)
        #expect(history.forward == [second])
        // From a place the history doesn't keep, such as the review, a visit only drops the
        // forward history.
        history.visit(from: nil)
        #expect(history.back == [.dashboard] && history.forward.isEmpty)
    }

    @Test func theBackHistoryKeepsTheLatestTwoHundredPlaces() {
        var history = VisitHistory()
        for i in 0..<205 { history.visit(from: .lesson(LessonRef(library: "library", lesson: "\(i)"))) }
        #expect(history.back.count == 200)
        #expect(history.back.first == .lesson(LessonRef(library: "library", lesson: "5")))
    }

    @Test func backAndForwardStepOverThePlaceOnScreen() {
        var history = VisitHistory()
        // Mission Control, the first lesson, the review, then the first lesson again.
        history.visit(from: .dashboard)
        history.visit(from: first)
        history.visit(from: nil)
        // The lesson's own entry is no step: Back goes to Mission Control.
        #expect(history.canGoBack(from: first))
        #expect(history.goBack(from: first) == .dashboard)
        #expect(history.back.isEmpty && history.forward == [first])
        #expect(!history.canGoBack(from: .dashboard) && history.canGoForward(from: .dashboard))
        #expect(!history.canGoForward(from: first))
        #expect(history.goForward(from: .dashboard) == first)
        #expect(history.back == [.dashboard] && history.forward.isEmpty)
        #expect(history.goForward(from: first) == nil)
    }

    @Test func backFromTheReviewLeavesNoForwardStepForIt() {
        var history = VisitHistory()
        history.visit(from: first)
        #expect(history.goBack(from: nil) == first)
        #expect(history.back.isEmpty && history.forward.isEmpty)
    }

    @Test func thePlaceOnScreenIsWhatTheHistoryKeeps() {
        let navigation = NavigationModel()
        #expect(navigation.visited == nil && navigation.lessonTarget == nil)
        navigation.screen = .dashboard
        #expect(navigation.visited == .dashboard)
        navigation.screen = .fieldGuide
        #expect(navigation.visited == .fieldGuide)
        // The documentation is no place before it shows a page; then it is where it is.
        navigation.screen = .documentation
        #expect(navigation.visited == nil)
        navigation.docLocation = { .symbol(package: "fmt", name: "Println") }
        #expect(navigation.visited == .doc(.symbol(package: "fmt", name: "Println")))
        #expect(navigation.lessonTarget == nil)
        let ref = LessonRef(library: "library", lesson: "first")
        navigation.screen = .lesson(ref)
        #expect(navigation.visited == .lesson(ref) && navigation.lessonTarget == ref)
    }

    @Test func backAndForwardGoFromThePlaceOnScreen() {
        let navigation = NavigationModel()
        let ref = LessonRef(library: "library", lesson: "first")
        navigation.screen = .dashboard
        #expect(!navigation.canGoBack && !navigation.canGoForward)
        navigation.visit()
        navigation.screen = .lesson(ref)
        navigation.visit()
        navigation.screen = .fieldGuide
        #expect(navigation.canGoBack && !navigation.canGoForward)
        #expect(navigation.goBack() == first)
        // AppModel shows the place Back returns to.
        navigation.screen = .lesson(ref)
        #expect(navigation.history.forward == [.fieldGuide])
        #expect(navigation.canGoBack && navigation.canGoForward)
        #expect(navigation.goForward() == .fieldGuide)
        navigation.screen = .fieldGuide
        #expect(navigation.history.back == [.dashboard, first] && navigation.history.forward.isEmpty)
        // On to the review, which the history doesn't keep: Back from it puts nothing in the
        // forward history.
        navigation.visit()
        navigation.screen = nil
        #expect(navigation.goBack() == .fieldGuide && navigation.history.forward.isEmpty)
    }
}
