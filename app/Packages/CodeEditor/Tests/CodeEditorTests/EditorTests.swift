import AppKit
import SwiftUI
import SyntaxCore
import Testing

@testable import CodeEditor

@Suite struct EditorStyleTests {
    @Test func tabsAreFourSpacesWide() {
        let style = EditorStyle(fontSize: 13)
        let space = (" " as NSString).size(withAttributes: [.font: style.font]).width
        #expect(abs(style.paragraphStyle.defaultTabInterval - 4 * space) < 0.001)
        #expect(style.paragraphStyle.tabStops.isEmpty)
        #expect(style.font.fontName == NSFont.monospacedSystemFont(ofSize: 13, weight: .regular).fontName)
        #expect(style.font.pointSize == 13)
    }

    @Test func everyKindKeepsTheBaseParagraphAndFontSize() {
        let style = EditorStyle(fontSize: 15)
        for kind in SyntaxKind.allCases {
            let attributes = style.attributes(for: kind)
            #expect(attributes[.paragraphStyle] as? NSParagraphStyle === style.paragraphStyle)
            #expect((attributes[.font] as? NSFont)?.pointSize == 15)
            #expect((attributes[.foregroundColor] as? NSColor) === style.theme.color(for: kind))
        }
    }
}

@Suite @MainActor struct DebouncerTests {
    @Test func firesOnceAfterTheLastSchedule() async throws {
        let debouncer = Debouncer(delay: .milliseconds(50))
        var fired: [Int] = []
        debouncer.schedule { fired.append(1) }
        debouncer.schedule { fired.append(2) }
        #expect(debouncer.isPending)
        try await Task.sleep(for: .milliseconds(400))
        #expect(fired == [2])
        #expect(!debouncer.isPending)
    }

    @Test func flushRunsThePendingActionNow() async throws {
        let debouncer = Debouncer(delay: .milliseconds(50))
        var fired = 0
        debouncer.schedule { fired += 1 }
        debouncer.flush()
        #expect(fired == 1)
        try await Task.sleep(for: .milliseconds(200))
        #expect(fired == 1)
        debouncer.flush()
        #expect(fired == 1)
    }

    @Test func cancelDropsThePendingAction() async throws {
        let debouncer = Debouncer(delay: .milliseconds(50))
        var fired = 0
        debouncer.schedule { fired += 1 }
        debouncer.cancel()
        try await Task.sleep(for: .milliseconds(200))
        #expect(fired == 0)
    }
}

/// The text view's own behavior, without a window.
@Suite @MainActor struct CodeTextViewTests {
    private func makeTextView(_ text: String, selection: Int) -> CodeTextView {
        let textView = CodeTextView(usingTextLayoutManager: true)
        textView.configure()
        textView.string = text
        textView.setSelectedRange(NSRange(location: selection, length: 0))
        return textView
    }

    /// Code copied from a readme carries the characters that lay it out there; a paste leaves
    /// them out, and keeps everything else.
    @Test func aPasteLeavesOutCharactersThatNeverBelongInCode() {
        #expect(PastedText.clean("\u{202F}o\u{2060}s\u{2060}.\u{2060}Open\u{202F}") == "os.Open")
        #expect(PastedText.clean("a\u{200B}b\u{FEFF}c d\u{00A0}e") == "abc d\u{00A0}e")
        #expect(!PastedText.needsCleaning("plain code"))
        let textView = makeTextView("x := ", selection: 5)
        let board = NSPasteboard.withUniqueName()
        defer { board.releaseGlobally() }
        board.declareTypes([.string], owner: nil)
        board.setString("\u{202F}o\u{2060}s\u{2060}.\u{2060}Getwd()\u{202F}", forType: .string)
        #expect(textView.readSelection(from: board, type: .string))
        #expect(textView.string == "x := os.Getwd()")
    }

    @Test func usesTextKit2WithCodeFriendlySettings() {
        let textView = makeTextView("x", selection: 0)
        #expect(textView.textLayoutManager != nil)
        #expect(!textView.isAutomaticQuoteSubstitutionEnabled)
        #expect(!textView.isAutomaticDashSubstitutionEnabled)
        #expect(!textView.isAutomaticTextReplacementEnabled)
        #expect(!textView.isAutomaticSpellingCorrectionEnabled)
        #expect(!textView.isContinuousSpellCheckingEnabled)
        #expect(!textView.isGrammarCheckingEnabled)
        #expect(!textView.smartInsertDeleteEnabled)
        #expect(!textView.isRichText)
        #expect(textView.allowsUndo)
        #expect(textView.usesFindBar)
        #expect(textView.isIncrementalSearchingEnabled)
    }

    @Test func returnIndentsAndSplitsBraces() {
        let textView = makeTextView("func main() {}", selection: 13)
        textView.insertNewline(nil)
        #expect(textView.string == "func main() {\n\t\n}")
        #expect(textView.selectedRange() == NSRange(location: 15, length: 0))
        #expect(textView.textLayoutManager != nil)
    }

    @Test func closingBraceOutdents() {
        let textView = makeTextView("if ok {\n\t\t", selection: 10)
        textView.insertText("}", replacementRange: NSRange(location: NSNotFound, length: 0))
        #expect(textView.string == "if ok {\n\t}")
    }

    @Test func tabInsertsATab() {
        let textView = makeTextView("x", selection: 0)
        textView.insertTab(nil)
        #expect(textView.string == "\tx")
    }

    @Test func tabIndentsASelectionThatSpansLines() {
        let textView = makeTextView("a\nb", selection: 0)
        textView.setSelectedRange(NSRange(location: 0, length: 3))
        textView.insertTab(nil)
        #expect(textView.string == "\ta\n\tb")
        #expect(textView.selectedRange() == NSRange(location: 0, length: 5))
        // A whole line with its line break is a selection that spans lines too.
        textView.setSelectedRange(NSRange(location: 0, length: 3))
        textView.insertTab(nil)
        #expect(textView.string == "\t\ta\n\tb")
        // Within one line, Tab replaces the selection as usual.
        textView.setSelectedRange(NSRange(location: 2, length: 1))
        textView.insertTab(nil)
        #expect(textView.string == "\t\t\t\n\tb")
    }

    @Test func shiftTabOutdentsTheLineOrTheSelectedLines() {
        let textView = makeTextView("\t\tx", selection: 3)
        textView.insertBacktab(nil)
        #expect(textView.string == "\tx")
        #expect(textView.selectedRange() == NSRange(location: 2, length: 0))
        textView.string = "\ta\n\tb"
        textView.setSelectedRange(NSRange(location: 0, length: 5))
        textView.insertBacktab(nil)
        #expect(textView.string == "a\nb")
        #expect(textView.selectedRange() == NSRange(location: 0, length: 3))
        textView.string = "\tx"
        textView.isEditable = false
        textView.insertBacktab(nil)
        #expect(textView.string == "\tx")
    }

    @Test func typedTextStartsPlain() {
        let textView = makeTextView("", selection: 0)
        textView.typingAttributes = [.foregroundColor: NSColor.red]
        #expect(textView.typingAttributes[.foregroundColor] as? NSColor === textView.style.theme.color(for: .plain))
        #expect(textView.typingAttributes[.paragraphStyle] as? NSParagraphStyle === textView.style.paragraphStyle)
    }

    @Test func readOnlyIgnoresTyping() {
        let textView = makeTextView("x", selection: 1)
        textView.isEditable = false
        textView.insertNewline(nil)
        #expect(textView.string == "x")
    }

    /// Option-click asks Quick Help about the word whose letters are under the pointer: never
    /// about the blank space past the end of a line, nor punctuation beside a word.
    @Test func optionClickFindsOnlyTheWordUnderThePointer() throws {
        let text = "fmt.Println(x)\n\nvar b strings.Builder\n"
        // The editor as the app has it, laid out in a window that is never shown.
        let coordinator = CodeEditorView(text: .constant(text), highlighter: TestHighlighter()).makeCoordinator()
        let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 600, height: 300), styleMask: [.titled], backing: .buffered, defer: true)
        window.contentView = coordinator.scrollView
        window.layoutIfNeeded()
        let textView = coordinator.textView
        let layout = try #require(textView.textLayoutManager)
        layout.ensureLayout(for: layout.documentRange)
        let ns = text as NSString
        func center(of index: Int) throws -> NSPoint {
            let rect = try #require(textView.rect(for: NSRange(location: index, length: 1)))
            return NSPoint(x: rect.midX, y: rect.midY)
        }
        let println = ns.range(of: "Println")
        #expect(textView.wordIndex(at: try center(of: println.location + 2)) == println.location + 2)
        // Its last letter's right half, where the insertion point would go after it.
        let last = try #require(textView.rect(for: NSRange(location: NSMaxRange(println) - 1, length: 1)))
        #expect(textView.wordIndex(at: NSPoint(x: last.maxX - 1, y: last.midY)) == NSMaxRange(println) - 1)
        // Punctuation, and the blank space past the end of a line ending in a word.
        #expect(textView.wordIndex(at: try center(of: ns.range(of: "(").location)) == nil)
        let builder = ns.range(of: "Builder")
        let end = try #require(textView.rect(for: NSRange(location: NSMaxRange(builder) - 1, length: 1)))
        #expect(textView.wordIndex(at: NSPoint(x: end.maxX + 120, y: end.midY)) == nil)
        // The empty line.
        let empty = try #require(textView.rect(for: NSRange(location: ns.range(of: "\n\n").location + 1, length: 1)))
        #expect(textView.wordIndex(at: NSPoint(x: empty.minX + 40, y: empty.midY)) == nil)
    }
}

/// What Option-click asked Quick Help about.
@MainActor private final class QuickHelpAsks {
    var offsets: [Int] = []
}

/// Events for the app to hand out in place of its own queue, as a text view tracking the mouse
/// reads them. The real queue is no use here: posting to it, and the window server's events
/// once the app reads them, stop the main run loop's current run. When that is the run the
/// tests run in, the test process ends at once, tests still to run, and `swift test` reports
/// no failure.
private final class EventQueue {
    private(set) var events: [NSEvent]

    init(_ events: [NSEvent]) {
        self.events = events
    }

    /// Runs `body` with `NSApplication.nextEvent(matching:until:inMode:dequeue:)` taking its
    /// events from this queue, and from nowhere else.
    @MainActor func serve(_ body: () -> Void) {
        _ = NSApplication.shared
        let selector = #selector(NSApplication.nextEvent(matching:until:inMode:dequeue:))
        guard let method = class_getInstanceMethod(NSApplication.self, selector) else {
            Issue.record("NSApplication has no nextEvent(matching:until:inMode:dequeue:)")
            return
        }
        let next: @convention(block) (NSApplication, UInt64, NSDate?, NSString?, Bool) -> NSEvent? = { _, mask, expiration, _, dequeue in
            self.next(matching: NSEvent.EventTypeMask(rawValue: mask), until: expiration as Date?, dequeue: dequeue)
        }
        let implementation = imp_implementationWithBlock(next)
        let original = method_setImplementation(method, implementation)
        defer {
            method_setImplementation(method, original)
            imp_removeBlock(implementation)
        }
        body()
    }

    private func next(matching mask: NSEvent.EventTypeMask, until expiration: Date?, dequeue: Bool) -> NSEvent? {
        guard let index = events.firstIndex(where: { mask.contains(NSEvent.EventTypeMask(type: $0.type)) }) else {
            // The app would wait for the person to do more.
            if let expiration, expiration > Date() {
                Issue.record("An event was waited for after the last one")
            }
            return nil
        }
        return dequeue ? events.remove(at: index) : events[index]
    }
}

/// Presses on the code as the app gets them, in a window that is never shown: `mouseDown`
/// gets the press, and asks the app for the moves and the release after it.
@Suite(.serialized) @MainActor struct MousePressTests {
    private let text = "fmt.Println(x)\nfmt.Println(y)\nfmt.Println(z)\n"
    private let lineLength = 15

    private func editor() throws -> (textView: CodeTextView, asks: QuickHelpAsks, window: NSWindow) {
        let asks = QuickHelpAsks()
        let coordinator = CodeEditorView(
            text: .constant(text), highlighter: TestHighlighter(), onOptionClick: { asks.offsets.append($0) }
        ).makeCoordinator()
        // Not deferred, so the window has a number for the events to carry. An event without
        // one has its location taken as on the screen when the text view tracks a click.
        let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 600, height: 300), styleMask: [.titled], backing: .buffered, defer: false)
        window.contentView = coordinator.scrollView
        window.layoutIfNeeded()
        let layout = try #require(coordinator.textView.textLayoutManager)
        layout.ensureLayout(for: layout.documentRange)
        return (coordinator.textView, asks, window)
    }

    /// A point a quarter of the way into the character at `index`, so the insertion point
    /// nearest it is the one before that character.
    private func point(at index: Int, in textView: CodeTextView) throws -> NSPoint {
        let rect = try #require(textView.rect(for: NSRange(location: index, length: 1)))
        return NSPoint(x: rect.minX + rect.width / 4, y: rect.midY)
    }

    private func mouse(_ type: NSEvent.EventType, at point: NSPoint, flags: NSEvent.ModifierFlags, in textView: NSTextView) throws -> NSEvent {
        try #require(NSEvent.mouseEvent(
            with: type, location: textView.convert(point, to: nil), modifierFlags: flags,
            timestamp: ProcessInfo.processInfo.systemUptime, windowNumber: textView.window?.windowNumber ?? 0,
            context: nil, eventNumber: 0, clickCount: 1, pressure: type == .leftMouseUp ? 0 : 1
        ))
    }

    /// Presses the button at `start` with `flags` held, moves through `path` and lets go at
    /// its end. Returns what the text view left in the queue, which the app would hand to the
    /// window next.
    private func press(_ textView: CodeTextView, at start: NSPoint, through path: [NSPoint] = [],
                       flags: NSEvent.ModifierFlags) throws -> [NSEvent.EventType] {
        let down = try mouse(.leftMouseDown, at: start, flags: flags, in: textView)
        let moves = try path.map { try mouse(.leftMouseDragged, at: $0, flags: flags, in: textView) }
        let queue = EventQueue(moves + [try mouse(.leftMouseUp, at: path.last ?? start, flags: flags, in: textView)])
        queue.serve { textView.mouseDown(with: down) }
        return queue.events.map(\.type)
    }

    private func selection(of textView: NSTextView) -> [NSRange] {
        textView.selectedRanges.map(\.rangeValue)
    }

    /// Option-click on a name asks Quick Help about it when the button comes up, and leaves
    /// the selection as it was. The release is the window's, as after any click.
    @Test func optionClickOnANameAsksQuickHelp() throws {
        let (textView, asks, _) = try editor()
        let index = (text as NSString).range(of: "Println(y)").location + 2
        let left = try press(textView, at: try point(at: index, in: textView), flags: .option)
        #expect(asks.offsets == [index])
        #expect(selection(of: textView) == [NSRange(location: 0, length: 0)])
        #expect(left == [.leftMouseUp])
    }

    /// A trackpad click can move the pointer a point or two before the button comes up; it is
    /// still a click, up to `clickWobble` points away.
    @Test func aClickThatWobblesIsStillAClick() throws {
        let (textView, asks, _) = try editor()
        let index = (text as NSString).range(of: "Println(z)").location + 3
        let start = try point(at: index, in: textView)
        let path = [
            NSPoint(x: start.x + 1, y: start.y), NSPoint(x: start.x + 1, y: start.y - 2),
            NSPoint(x: start.x - CodeTextView.clickWobble, y: start.y),
        ]
        let left = try press(textView, at: start, through: path, flags: .option)
        #expect(asks.offsets == [index])
        #expect(selection(of: textView) == [NSRange(location: 0, length: 0)])
        // The moves were dropped; the release is still the window's.
        #expect(left == [.leftMouseUp])
    }

    /// An Option-drag selects a column of text, whether it starts on a name or anywhere else,
    /// and asks nothing. The editor stays on TextKit 2.
    @Test func anOptionDragSelectsAColumn() throws {
        let (textView, asks, _) = try editor()
        // From the P of the first line's Println to the second n of the third's: a wobble,
        // then on to the end.
        let start = try point(at: 4, in: textView)
        let end = try point(at: 2 * lineLength + 10, in: textView)
        let halfway = NSPoint(x: (start.x + end.x) / 2, y: (start.y + end.y) / 2)
        var left = try press(textView, at: start, through: [NSPoint(x: start.x + 1, y: start.y + 1), halfway, end], flags: .option)
        #expect(selection(of: textView) == [
            NSRange(location: 4, length: 6), NSRange(location: lineLength + 4, length: 6),
            NSRange(location: 2 * lineLength + 4, length: 6),
        ])
        #expect(left.isEmpty)

        // From the first line's period, which is no name, to the second line's l.
        left = try press(textView, at: try point(at: 3, in: textView), through: [try point(at: lineLength + 9, in: textView)], flags: .option)
        #expect(selection(of: textView) == [NSRange(location: 3, length: 6), NSRange(location: lineLength + 3, length: 6)])
        #expect(left.isEmpty)
        #expect(asks.offsets.isEmpty)
        #expect(textView.textLayoutManager != nil)
    }

    /// Off a name, Option-click is an ordinary click, as is a click without Option: the text
    /// view puts the insertion point there.
    @Test func otherClicksAreTheTextViewsOwn() throws {
        let (textView, asks, _) = try editor()
        let ns = text as NSString
        let parenthesis = ns.range(of: "(y)").location
        var left = try press(textView, at: try point(at: parenthesis, in: textView), flags: .option)
        #expect(selection(of: textView) == [NSRange(location: parenthesis, length: 0)])
        #expect(left.isEmpty)

        let name = ns.range(of: "Println(z)").location + 1
        left = try press(textView, at: try point(at: name, in: textView), flags: [])
        #expect(selection(of: textView) == [NSRange(location: name, length: 0)])
        #expect(left.isEmpty)
        #expect(asks.offsets.isEmpty)
    }
}

/// The SwiftUI view's update path, driven through its coordinator in a window that is never shown.
@Suite @MainActor struct CodeEditorViewTests {
    private let code = (1...300).map { "line\($0) := \($0) // note" }.joined(separator: "\n")

    private func editor(showsLineNumbers: Bool) -> CodeEditorView {
        CodeEditorView(text: .constant(code), highlighter: TestHighlighter(), showsLineNumbers: showsLineNumbers)
    }

    @Test func lineNumbersToggleLiveAndKeepScrollAndSelection() throws {
        let coordinator = editor(showsLineNumbers: true).makeCoordinator()
        let scrollView = coordinator.scrollView
        let window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 600, height: 400), styleMask: [.titled], backing: .buffered, defer: true
        )
        window.contentView = scrollView
        window.layoutIfNeeded()
        #expect(scrollView.rulersVisible)
        #expect(scrollView.verticalRulerView?.ruleThickness ?? 0 > 0)

        coordinator.textView.setSelectedRange(NSRange(location: 1_234, length: 5))
        let layoutManager = try #require(coordinator.textView.textLayoutManager)
        layoutManager.ensureLayout(for: layoutManager.documentRange)
        window.layoutIfNeeded()
        let clipView = scrollView.contentView
        clipView.scroll(to: NSPoint(x: clipView.bounds.minX, y: 800))
        scrollView.reflectScrolledClipView(clipView)
        #expect(clipView.bounds.minY == 800)
        let widthWithNumbers = coordinator.textView.frame.width

        coordinator.update(editor(showsLineNumbers: false))
        window.layoutIfNeeded()
        #expect(!scrollView.rulersVisible)
        #expect(clipView.bounds.minY == 800)
        #expect(coordinator.textView.selectedRange() == NSRange(location: 1_234, length: 5))
        #expect(coordinator.textView.frame.width > widthWithNumbers)

        coordinator.update(editor(showsLineNumbers: true))
        window.layoutIfNeeded()
        #expect(scrollView.rulersVisible)
        #expect(clipView.bounds.minY == 800)
        #expect(coordinator.textView.selectedRange() == NSRange(location: 1_234, length: 5))
        #expect(coordinator.textView.frame.width == widthWithNumbers)
        #expect(coordinator.textView.textLayoutManager != nil)
    }
}

/// The code an editor writes back through its binding, and how often it did.
@MainActor private final class TextBox {
    var text: String
    var writes = 0

    init(_ text: String) {
        self.text = text
    }

    var binding: Binding<String> {
        Binding(get: { self.text }, set: { self.text = $0; self.writes += 1 })
    }
}

/// Commands from outside the editor, through its proxy, in a window that is never shown.
@Suite @MainActor struct CodeEditorProxyTests {
    private func editor(_ box: TextBox, fileName: String = "main.go", isEditable: Bool = true,
                        proxy: CodeEditorProxy) -> (CodeEditorView.Coordinator, NSWindow) {
        let coordinator = CodeEditorView(
            text: box.binding, highlighter: TestHighlighter(), lineSyntax: LineSyntax(fileName: fileName),
            isEditable: isEditable, proxy: proxy
        ).makeCoordinator()
        let window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 600, height: 400), styleMask: [.titled], backing: .buffered, defer: true
        )
        window.contentView = coordinator.scrollView
        return (coordinator, window)
    }

    /// Ends the event, as the app's run loop does after each one: the undo manager closes the
    /// group it opened for the event's changes.
    private func endEvent() {
        RunLoop.current.run(mode: .default, before: Date())
    }

    @Test func eachCommandIsOneUndoStep() throws {
        let proxy = CodeEditorProxy()
        let box = TextBox("func f() {\nx := 1\n}")
        let (coordinator, window) = editor(box, proxy: proxy)
        let textView = coordinator.textView
        textView.setSelectedRange(NSRange(location: 11, length: 0))
        // No event ends between the commands: each closes its own undo step.
        #expect(proxy.indent())
        #expect(textView.string == "func f() {\n\tx := 1\n}" && box.text == textView.string)
        #expect(proxy.toggleComment())
        #expect(textView.string == "func f() {\n\t// x := 1\n}" && box.text == textView.string)
        #expect(textView.selectedRange() == NSRange(location: 15, length: 0))
        #expect(window.firstResponder === textView)

        let undo = try #require(textView.undoManager)
        #expect(undo.undoActionName == "Toggle Comment")
        undo.undo()
        #expect(textView.string == "func f() {\n\tx := 1\n}" && box.text == textView.string)
        #expect(undo.undoActionName == "Indent")
        undo.undo()
        #expect(textView.string == "func f() {\nx := 1\n}" && box.text == textView.string)
        #expect(!undo.canUndo)
        undo.redo()
        #expect(textView.string == "func f() {\n\tx := 1\n}" && box.text == textView.string)
        #expect(proxy.outdent())
        #expect(textView.string == "func f() {\nx := 1\n}")
    }

    @Test func replacingTheCodeKeepsTheCaretAndUndoesInOneStep() throws {
        let proxy = CodeEditorProxy()
        let original = "package main\n\nfunc main() {\nx:=1\n_ = x\n}\n"
        let formatted = "package main\n\nfunc main() {\n\tx := 1\n\t_ = x\n}\n"
        let box = TextBox(original)
        let (coordinator, _) = editor(box, proxy: proxy)
        let textView = coordinator.textView
        // Line 5, "_ = x", column 2.
        textView.setSelectedRange(NSRange(location: (original as NSString).range(of: "_ = x").location + 2, length: 0))
        #expect(proxy.replaceText(with: formatted, actionName: "Format"))
        #expect(textView.string == formatted && box.text == formatted && proxy.text == formatted)
        let caret = (formatted as NSString).range(of: "\t_ = x").location + 2
        #expect(textView.selectedRange() == NSRange(location: caret, length: 0))
        // gofmt's output lands after the key press's event has ended; typing that follows it
        // in the next event is a step of its own.
        textView.insertText("y", replacementRange: NSRange(location: NSNotFound, length: 0))
        endEvent()

        let undo = try #require(textView.undoManager)
        #expect(undo.undoActionName == "Typing")
        undo.undo()
        #expect(textView.string == formatted && box.text == formatted)
        #expect(undo.undoActionName == "Format")
        undo.undo()
        #expect(textView.string == original && box.text == original && proxy.text == original)
        #expect(!undo.canUndo)
    }

    /// TextKit 2 text views say nothing to their delegate when undo changes the text; the
    /// binding still gets the text back.
    @Test func undoingTypingWritesTheTextBack() throws {
        let proxy = CodeEditorProxy()
        let box = TextBox("x")
        let (coordinator, _) = editor(box, proxy: proxy)
        let textView = coordinator.textView
        textView.setSelectedRange(NSRange(location: 1, length: 0))
        textView.insertText("y", replacementRange: NSRange(location: NSNotFound, length: 0))
        #expect(box.text == "xy")
        endEvent()
        let undo = try #require(textView.undoManager)
        undo.undo()
        #expect(textView.string == "x" && box.text == "x")
        undo.redo()
        #expect(textView.string == "xy" && box.text == "xy")
    }

    @Test func theSameCodeChangesNothing() {
        let proxy = CodeEditorProxy()
        let box = TextBox("package main\n")
        let (coordinator, _) = editor(box, proxy: proxy)
        #expect(!proxy.replaceText(with: "package main\n", actionName: "Format"))
        #expect(box.writes == 0)
        #expect(coordinator.textView.undoManager?.canUndo == false)
    }

    @Test func aReadOnlyEditorTakesNoCommands() {
        let proxy = CodeEditorProxy()
        let box = TextBox("\tx := 1")
        let (coordinator, _) = editor(box, isEditable: false, proxy: proxy)
        #expect(!proxy.indent() && !proxy.outdent() && !proxy.toggleComment())
        #expect(!proxy.replaceText(with: "x := 1", actionName: "Format"))
        #expect(coordinator.textView.string == "\tx := 1" && box.writes == 0)
    }

    @Test func toggleCommentNeedsAFileWithLineComments() {
        let proxy = CodeEditorProxy()
        let box = TextBox("x")
        _ = editor(box, fileName: "notes.txt", proxy: proxy)
        #expect(!proxy.toggleComment())
        #expect(box.text == "x" && box.writes == 0)
    }

    @Test func commandsShowHighlightedAtOnce() throws {
        let proxy = CodeEditorProxy()
        let box = TextBox("x := 1")
        let (coordinator, _) = editor(box, proxy: proxy)
        let storage = try #require(coordinator.textView.textStorage)
        #expect(proxy.toggleComment())
        // Painted as a comment already, not a moment later when highlighting runs.
        let color = storage.attribute(.foregroundColor, at: 4, effectiveRange: nil) as? NSColor
        #expect(color === coordinator.textView.style.theme.color(for: .comment))
    }

    @Test func theProxyActsOnTheNewestEditor() {
        let proxy = CodeEditorProxy()
        #expect(proxy.text == nil && !proxy.indent())
        let first = CodeEditorView(text: .constant("first"), highlighter: TestHighlighter(), proxy: proxy).makeCoordinator()
        #expect(proxy.text == "first")
        let second = CodeEditorView(text: .constant("second"), highlighter: TestHighlighter(), proxy: proxy).makeCoordinator()
        #expect(proxy.text == "second")
        // The first editor going away leaves the second one attached.
        first.dismantle()
        #expect(proxy.text == "second")
        second.dismantle()
        #expect(proxy.text == nil)
    }
}
