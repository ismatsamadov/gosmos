# Design

Gosmos is a native learning space with a quiet space-explorer identity: native controls,
Liquid Glass actions, clear system typography, and every action on the keyboard. The rounded
wordmark and a small native orbit illustration give it an identity without crowding the work. It
never uses the network: every lesson, code or quiz, is checked on the Mac.

## Mission Control and rewards

- The app opens to Mission Control. Continue Learning resumes an unfinished lesson or picks
  the first unfinished one. Course cards show required-lesson progress and open the next
  unfinished lesson in that course. View › Mission Control (⌘1) returns here.
- The dashboard leads with the next lesson's title, course, chapter and type on an ink-colored
  panel. A single summary strip shows XP, completed lessons, rank and reviews; the review
  summary opens Review. Below, the learning path sits beside compact achievement badges and
  a documentation shortcut. Narrow panes stack the two columns.
- With a single course, its card shows three chapters near the next unfinished lesson.
  Each chapter opens its next unfinished lesson; completed chapters can be revisited. The
  course action always opens the next lesson in that course. An Add Library action sits
  beside the learning-path heading.
- Each completed lesson contributes 100 XP; every 500 XP advances a level. Ranks begin at
  Cadet, then Explorer, Navigator, Pilot, Commander, and Captain; levels keep advancing.
- Six badges recognize a first lesson, coding, a quiz, five lessons, a chapter, and 25
  lessons. A chapter badge requires all of its required lessons, with at least one required
  lesson. Optional lessons earn XP too.
- XP, rank and badges are computed from saved completions in the open libraries. Existing
  progress counts immediately. Reopening or resubmitting a lesson never adds duplicate XP;
  marking it not done removes its contribution. No new database or reward ledger is needed.
- The sidebar keeps rank, XP and the next level visible. Lessons show chapter progress and
  their XP value, and a pass shows a completion panel with a next-mission action. A review
  shows session progress and encouragement, without awarding repeat completion XP.
- Motivation comes from visible progress, achievable milestones, and encouragement. There
  are no streak penalties, competitive leaderboards, or timers.

## Platform

- macOS 27 or later, the current release, so the app is Liquid Glass throughout. The standard containers get it
  from the system: the split view's sidebar, the unified toolbar, sheets, menus and the
  inspector. Glass belongs to actions and system chrome; reading, editor, console and quiz
  content use quiet solid surfaces. There is no hand-made blur, gradient or shadow.
- Light, dark, Increase Contrast and Reduce Motion all look right.
- Text uses semantic labels and contrast-aware colors. Go Blue, #00ADD8 in dark mode and
  Go Teal #00758D in light mode, is the primary accent. Gold identifies XP and achievements;
  violet identifies rank, review and the documentation field guide. A pale gray canvas in
  light mode and charcoal canvas in dark mode, solid surfaces and fine borders define cards.
  The dashboard's welcome panel stays dark teal in either appearance. Headings use semibold
  system type with restrained spacing; rounded type is reserved for the wordmark.
- Icons are SF Symbols with hierarchical rendering. State changes use symbol effects:
  `.bounce` on success, `.replace` between Run and Stop.

## Window

- One main window, at least 1100×700. It remembers its size, the sidebar width, every split
  and the last lesson. The sidebar starts 260 pt wide and the inspector 240, so with the
  lesson's 600 they fit the narrowest window.
- `NavigationSplitView` with the sidebar and the lesson. The toolbar:
  - back and forward through the lessons and documentation pages you visited;
  - the lesson's title, with course › chapter as the subtitle;
  - Run ⌘R, and Submit ⌘↩ as the prominent glass button;
  - the search field: Edit › Search (⇧⌘F, as Xcode's Find in Project) puts the insertion
    point in it, or in the documentation's own field while the documentation is on screen;
    ⌘F stays the editor's Find;
  - the inspector toggle.
- The inspector on the right holds your note for the lesson (autosaved) and the lesson's
  facts: type, difficulty, and topics as tags.

## Sidebar

- Gosmos branding, then Mission Control, Review and Go Documentation (`text.book.closed`),
  then the libraries; a compact rank and progress card sits at the bottom. The sidebar's
  background is consistent from the wordmark through its list and footer.
- One section per course, with its title and a progress ring for finished lessons out of
  required ones. Guided projects sit under their own heading.
- Chapters are disclosure groups; the current one is open.
- Each lesson row shows a type symbol, the title, and a done mark (`checkmark.circle.fill`
  in the tint). A title too long for one line takes a second, breaking between words, never
  inside a hyphenated one; chapter titles too. The type symbols:
  - code: `chevron.left.forwardslash.chevron.right`
  - tests: `checklist`
  - quiz: `questionmark.bubble`
- Search runs through the full-text index and lists matching lessons, with the match
  highlighted, then a Documentation section with the best matches in the Go documentation.

## Reading

- The readme column is 680 pt, room for 80 columns of code, centered in its pane. Body text
  keeps to 560 pt of it, about 80 characters a line; code blocks, tables and panels take the
  whole column. Body text is SF Pro 15 pt with a 1.45 line height, and paragraphs are 12 pt
  apart. Headings are 28, 22, 18 and 15 pt semibold. All text is selectable.
- Text scrolled up under the toolbar stops at a hard scroll edge, so it never runs into the
  window's title.
- Code blocks are SF Mono 13 pt on a 10 pt rounded `.quaternary` panel, with a language
  label and a Copy button on hover. Go is highlighted in Xcode's default light and dark
  palettes.
- Callouts ([!NOTE], [!TIP], [!IMPORTANT], [!WARNING], [!CAUTION]) are tinted rounded
  panels with their symbol.
- Tables have hairline rules and zebra rows. Inline code sits in a subtle capsule. A link to
  a lesson of the library or to the Go documentation opens it in the app; any other address
  is plain text, since Gosmos opens nothing outside itself.
- Inline code that names a standard library package or a name in one (`fmt.Println`,
  `*strings.Builder`, `net/http`) links to its documentation, in the tint; Back returns to
  the lesson. Code blocks never link.

## Lessons

- **Code**:
  - A compact, open lesson heading shows chapter position, title, type, XP and chapter
    progress. It shares the reading surface instead of taking a separate card.
  - Readme on the left. On the right, the editor above the console.
  - The editor is SF Mono 13 pt (size in Settings), with a quiet line-number gutter, a
    current-line highlight, tabs four wide, highlighting and the find bar.
  - The file bar above it ends with borderless buttons, also in the Editor menu: Format
    (`wand.and.stars`, gofmt, ⇧⌥F), Toggle Comment (`//`, ⌘/), Outdent (⌥⌘[) and Indent
    (⌥⌘]): 14 pt glyphs in the label color, 28×24 pt, with a rounded highlight under the
    pointer and a deeper one while pressed, as Xcode's editor buttons. Each change is one
    undo step. A button that can't act, as on a read-only file, says
    why in its help. Code gofmt can't parse stays as it is, and the console shows gofmt's
    messages under Format.
  - Quick Help: ⌥-click a package's name or a name it qualifies, such as `strings.Builder` or
    `fmt.Println`, and a popover shows its declaration, the first paragraph of its doc and
    Open in Documentation (Return). The file's imports say which package a name is in,
    renamed ones too. Editor › Show Documentation (⌃⌘?, Xcode's Quick Help keys) does the
    same for the name at the insertion point. A name outside the standard library gets the
    system beep.
  - The console opens with a compact status bar: state symbol, duration, passed or failed.
    The output below it is monospaced, and the first failing test is expanded.
  - "Reset to Starter Code" asks first.
  - Under the readme, hints show one at a time and stay shown; the last opens after a Submit
    that fails. A Show Solution row (⌥⌘S) says that looking before you pass brings the lesson
    back for review. It opens a sheet, 1040 pt wide, with your code and the solution side by
    side, the lines that differ marked, and why the solution works. Nothing goes into the editor. A long file
    is compared in the background, with a small spinner meanwhile.
- **Quiz**:
  - Each question is a card, and its answers are large selectable rows. Keys 1 to 5 pick an
    answer. Number keys, radio marks, a tinted selection and a fine outline make each state
    clear. The answered count and a slim progress track sit above the questions.
  - After Submit the right answer turns green with a checkmark. A wrong pick turns red, and
    the right one is shown. Under every answer, a line says why it is right or wrong.
  - Success gives a light trackpad haptic (NSHapticFeedbackManager) and the done mark
    bounces.
- After a pass, Next Lesson (⌘]) appears in the toolbar. A pass that schedules a review says
  so in the console.
- **Review**: the sidebar's Review entry has a badge with the number of items due. Due quiz
  questions come one at a time, lessons taking turns, then the lessons to redo, each with
  its reason. A redo starts from the starter code in a folder of its own and isn't kept.
  The session shows how many questions were practiced and a progress bar, with a reassuring
  completion state. No streaks, counts of misses, or reminders.

## Documentation

The Go documentation is built into the app, so looking something up never needs the
internet: The Go Programming Language Specification, The Go Memory Model, Effective Go, the
FAQ, How to Write Go Code, and every package of the standard library, for the Go built in.

- Go Documentation in the sidebar shows the browser in the detail area: a list beside the
  page, with a split that is remembered. The window's title is the page's, with Go
  Documentation as the subtitle.
- The list has Language, the five articles, then Standard Library: the packages grouped by
  the first element of their path (archive, crypto, encoding…), the group of the page on
  screen open. Acknowledgements ends it. Arrow keys move through it, opening each page.
- A search field heads the list (⌘F goes to it, Escape clears it, Return opens the first
  result), and stays there when nothing is found. Results
  come as you type, kind symbol, qualified name and synopsis: the exact name first, so
  `Println` puts fmt.Println on top and `Builder` strings.Builder, then names that start
  with what you typed, then the full text of the docs and the articles' sections.
- Pages use the readme's column and type. A package's page: Package, its path and import
  line, its doc, its examples, an Index of one-line signatures that jump to their
  declarations, then Constants, Variables, Functions and Types, each type followed by its
  funcs and methods. A declaration is Go highlighted as in the editor, then its doc, then
  its examples: the code, and under Output what it prints. An article's page: its title and
  version, Contents (the sections two levels down, each a link), then the sections.
- Deprecated packages and names are marked quietly: a small gray Deprecated capsule, as the
  operating-system tags of headings look; a deprecated package is gray in the list.
- Links in the documentation (`godoc:`, `godoc-article:`) open in place, scrolled to the
  name or section. Back and Forward (⌥⌘← / ⌥⌘→) go through documentation pages and lessons
  alike.
- Pages lay out lazily, so the specification or net/http opens as fast as a small package.
- The Acknowledgements page, also in Settings › General, shows the documentation's credit:
  Go's BSD-3-Clause license, and the CC BY 4.0 credit to The Go Authors for the articles.

## States

- Use `ContentUnavailableView` when no library is open (with Open Library…), a library is
  refused (naming the problems), nothing is selected, a search finds nothing, or Go can't run
  (saying why, with a button to Settings).
- Errors appear as an inline banner or an alert with a clear action. They are never silent.
- Long work shows progress in place, never a frozen window. A spinner shows only once a wait
  has lasted a moment (0.4 s), so a quick load never flashes one; the page or lesson on
  screen stays until the next one is ready.

## Settings

A standard Settings window with grouped forms:
- General: the Go documentation's Acknowledgements…; the libraries, the built-in courses
  first and then each file you opened, with their lesson counts and how many you have done;
  Open Library….
- Editor: font size, line numbers.
- Go: which Go runs your code: "Built-in Go 1.27.3" unless you chose another, then its
  version and path. When it can't run, why, in place of the version. Below, as an option for
  advanced use, Choose… picks another Go on this Mac, and Use Built-in Go goes back. Nothing
  links to a download or a web page: the app works offline.

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
