# Design

Gosmos should look and feel like an app Apple would ship: native controls, Liquid Glass,
calm typography, and every action on the keyboard. Nothing in it looks like a web page.

## Platform

- macOS 27 or later, the current release, so the app is Liquid Glass throughout. The standard containers get it
  from the system: the split view's sidebar, the unified toolbar, sheets, menus and the
  inspector. Custom surfaces use `glassEffect` inside a `GlassEffectContainer`: the run
  status bar, quiz answers and command-line steps. There is no hand-made blur, gradient or
  shadow.
- Light, dark, Increase Contrast and Reduce Motion all look right.
- Colors are semantic: `.primary`, `.secondary`, `.tint` and system backgrounds. The one
  accent is Go Blue, #00ADD8 in dark mode and Go Teal #00758D in light mode for contrast,
  set as the asset's AccentColor.
- Icons are SF Symbols with hierarchical rendering. State changes use symbol effects:
  `.bounce` on success, `.replace` between Run and Stop.

## Window

- One main window, at least 1100×700. It remembers its size, the sidebar width, every split
  and the last lesson.
- `NavigationSplitView` with the sidebar and the lesson. The toolbar:
  - back and forward through the lessons you visited;
  - the lesson's title, with course › chapter as the subtitle;
  - Run ⌘R, and Submit ⌘↩ as the prominent glass button;
  - the search field;
  - the inspector toggle.
- The inspector on the right holds your note for the lesson (autosaved) and the lesson's
  facts: type, difficulty, and topics as tags.

## Sidebar

- One section per course, with its title and a progress ring for finished lessons out of
  required ones. Guided projects sit under their own heading.
- Chapters are disclosure groups; the current one is open.
- Each lesson row shows a type symbol, the title, and a done mark (`checkmark.circle.fill`
  in the tint). The type symbols:
  - code: `chevron.left.forwardslash.chevron.right`
  - tests: `checklist`
  - quiz: `questionmark.bubble`
  - command line: `terminal`
  - GitHub: `arrow.triangle.branch`
  - interview: `person.wave.2`
- Search runs through the full-text index and lists matching lessons, with the match
  highlighted.

## Reading

- The readme column is about 72 characters wide (about 680 pt), centered in its pane. Body
  text is SF Pro 15 pt with a 1.45 line height, and paragraphs are 12 pt apart. Headings
  are 28, 22, 18 and 15 pt semibold. All text is selectable.
- Code blocks are SF Mono 13 pt on a 10 pt rounded `.quaternary` panel, with a language
  label and a Copy button on hover. Go is highlighted in Xcode's default light and dark
  palettes.
- Callouts ([!NOTE], [!TIP], [!IMPORTANT], [!WARNING], [!CAUTION]) are tinted rounded
  panels with their symbol.
- Tables have hairline rules and zebra rows. Inline code sits in a subtle capsule. Links
  open in the default browser.

## Lessons

- **Code**:
  - Readme on the left. On the right, the editor above the console.
  - The editor is SF Mono 13 pt (size in Settings), with a quiet line-number gutter, a
    current-line highlight, tabs four wide, highlighting and the find bar.
  - The console opens with a glass status bar: state symbol, duration, passed or failed.
    The output below it is monospaced, and the first failing test is expanded.
  - "Reset to Starter Code" asks first.
  - Under the readme, hints show one at a time and stay shown; the last opens after a Submit
    that fails. A Show Solution row (⌥⌘S) says that looking before you pass brings the lesson
    back for review. It opens a sheet with your code and the solution side by side, the lines
    that differ marked, and why the solution works. Nothing goes into the editor. A long file
    is compared in the background, with a small spinner meanwhile.
- **Quiz**:
  - Each question is a card, and its answers are large selectable rows. Keys 1 to 5 pick an
    answer.
  - After Submit the right answer turns green with a checkmark. A wrong pick turns red, and
    the right one is shown. Under every answer, a line says why it is right or wrong.
  - Success gives a light trackpad haptic (NSHapticFeedbackManager) and the done mark
    bounces.
- **Command line**:
  - A pipeline of steps. Each step shows a state (pending, running, passed, failed), its
    command or `METHOD url`, and a detail that opens to show output, the response's status,
    headers and body, and each test passed or failed.
  - The project folder and server address sit at the top and can be edited.
- **GitHub and interview**: the readme, the checks as a reference checklist, and a clear
  Mark as Done button.
- After a pass, Next Lesson (⌘]) appears in the toolbar. A pass that schedules a review says
  so in the console.
- **Review**: the sidebar's Review entry has a badge with the number of items due. Due quiz
  questions come one at a time, lessons taking turns, then the lessons to redo, each with
  its reason. A redo starts from the starter code in a folder of its own and isn't kept.
  Calm on purpose: no streaks, no counts of misses, no reminders.

## States

- Use `ContentUnavailableView` when no library is open (with Open Library…), a library is
  refused (naming the problems), nothing is selected, a search finds nothing, or Go is
  missing (with a button to Settings).
- Errors appear as an inline banner or an alert with a clear action. They are never silent.
- Long work shows progress in place, never a frozen window.

## Settings

A standard Settings window with grouped forms:
- General: the libraries, the built-in courses first and then each file you opened, with
  their lesson counts and how many you have done; Open Library….
- Editor: font size, line numbers.
- Go: toolchain path and the detected version.

## Icon

A layered macOS 27 rounded-square icon: a Go Blue planet inside a tilted orbit, with a small
satellite on it, on a night-blue sky with a few stars. No text. Each size from 16 to 1024 px
is drawn on its own pixel grid, simpler at 16 and 32 px, so it stays sharp; delivered as
`AppIcon.icns` by `app/scripts/make-icon.sh`.

## Accessibility

- Every control has a VoiceOver label.
- The whole app works from the keyboard: Tab, arrow keys in the sidebar, and the shortcuts
  above.
- Hit targets are at least 24 pt.
