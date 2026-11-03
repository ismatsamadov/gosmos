package authoring

// The story in testdata/story is invented: the opening, a chapter's briefing and debrief, and
// the field guide of Learn Launches, the fictional course of the lessons in testdata/packages.

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/jobloom/gosmos/db"
)

func storyExample() string { return filepath.Join("testdata", "story") }

// exampleStory is what testdata/story holds.
var exampleStory = Story{
	Pages: []Page{
		{
			Title: "Welcome to Apogee", Kind: "text", Eyebrow: "LAUNCH SITE · MONDAY 1 MARCH 2027", Action: "Meet the team",
			Body: "Apogee builds small rockets, and the code that flies them.\n\n" +
				"Today you join the *flight software* team, three weeks before\nthe first launch of **Apogee One**.",
			Cards: []Card{
				{Title: "The Flight Director", Label: "FD", Subtitle: "Runs every countdown",
					Body: "Reads out the launch window, and asks each console for a `go` or a `no-go`."},
				{Title: "Pad 2", Label: "P2"},
			},
		},
		{
			Title: "Your Badge", Kind: "badge", Eyebrow: "APOGEE", Action: "Start the countdown",
			Body: "Everyone on the console wears one.",
			Cards: []Card{
				{Title: "Flight Software Engineer", Label: "Role"},
				{Title: "Console 4", Label: "Seat", Subtitle: "Next to the flight director"},
			},
		},
	},
	Chapters: []ChapterStory{{
		Chapter: "1-countdown",
		Briefing: "The first launch is three weeks away, and the countdown software prints nothing yet.\n\n" +
			"Today you make it read out the launch window.",
		Debrief: "The console reads the window out now, on time,\nevery time.",
	}},
	Terms: []Term{
		{Term: "go/no-go", Meaning: "The poll before a launch: each console says go, or calls a hold.", Linked: true},
		{Term: "launch window", Meaning: "The minutes in which a rocket can lift off and still reach its orbit.", Linked: true},
		{Term: "pad", Meaning: "The concrete a rocket stands on before it flies."},
		{Term: "T-minus", Meaning: "The time left before liftoff, counted down.", Linked: true},
	},
}

// storyDir makes a story package of the files given, by name.
func storyDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "story")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestReadStory(t *testing.T) {
	s, err := ReadStory(storyExample())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, exampleStory) {
		t.Errorf("ReadStory = %+v\nwant %+v", s, exampleStory)
	}
	if got := s.Summary(); got != "2 pages, 4 cards, 1 chapter and 4 terms" {
		t.Errorf("Summary = %q", got)
	}
	if got := (Story{Terms: s.Terms[:1]}).Summary(); got != "1 term" {
		t.Errorf("Summary of one term = %q", got)
	}
}

func TestReadWriteStoryRoundTrip(t *testing.T) {
	s, err := ReadStory(storyExample())
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "story")
	if err := s.Write(out); err != nil {
		t.Fatal(err)
	}
	sameFiles(t, out, storyExample())
}

// A folder with none of the files, dot files aside, is the story of a course without one,
// and that story is written as an empty folder.
func TestEmptyStory(t *testing.T) {
	s, err := ReadStory(storyDir(t, map[string]string{".DS_Store": "x"}))
	if err != nil || !s.Empty() || s.Summary() != "no story" {
		t.Fatalf("ReadStory of an empty folder = %+v, %v", s, err)
	}
	out := filepath.Join(t.TempDir(), "story")
	if err := s.Write(out); err != nil {
		t.Fatal(err)
	}
	if entries, err := os.ReadDir(out); err != nil || len(entries) != 0 {
		t.Errorf("the empty story wrote %v, %v; want an empty folder", entries, err)
	}
	if err := s.Write(storyDir(t, map[string]string{"notes.md": "to do\n"})); err == nil || !strings.Contains(err.Error(), "is not empty") {
		t.Errorf("Write into a folder that has files: err = %v", err)
	}
}

// Changed names the files in which two stories differ; a part left out is the same as an
// empty one.
func TestStoryChanged(t *testing.T) {
	s, _ := ReadStory(storyExample())
	other, _ := ReadStory(storyExample())
	if c := s.Changed(other); len(c) != 0 {
		t.Errorf("the same story changed %q", c)
	}
	other.Pages[1].Cards[0].Label = "Job"
	other.Terms = other.Terms[1:]
	if c := s.Changed(other); !slices.Equal(c, []string{"opening.md", "terms.md"}) {
		t.Errorf("Changed = %q, want opening.md and terms.md", c)
	}
	if c := s.Changed(Story{}); !slices.Equal(c, storyFiles) {
		t.Errorf("Changed against no story = %q", c)
	}
	if c := (Story{Pages: []Page{}}).Changed(Story{}); len(c) != 0 {
		t.Errorf("an empty opening and none differ in %q", c)
	}
}

// Into a library and back: Put, then LoadStory, gives the same story and the same files. Put
// replaces the course's whole story and leaves other courses' alone.
func TestPutAndLoadStory(t *testing.T) {
	d, err := db.Create(filepath.Join(t.TempDir(), "apogee.sqlite"), db.Info{Title: "Apogee", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, s := range []string{
		`INSERT INTO course VALUES ('c', 'learn-launches', 'Learn Launches', 'Course', 'Launch code.', 'An invented course.', 0)`,
		`INSERT INTO chapter VALUES ('ch', 'c', 0, '1-countdown', 'Countdown', '')`,
		`INSERT INTO course VALUES ('c2', 'learn-landings', 'Learn Landings', 'Course', 'Landing code.', 'Another invented course.', 1)`,
		`INSERT INTO chapter VALUES ('ch2', 'c2', 0, '1-countdown', 'Countdown', '')`,
	} {
		if _, err := d.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	put := func(s Story, course string) error {
		tx, err := d.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		if err := s.Put(tx, course); err != nil {
			return err
		}
		return tx.Commit()
	}
	load := func(course string) Story {
		s, err := LoadStory(d, course)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	s, err := ReadStory(storyExample())
	if err != nil {
		t.Fatal(err)
	}
	for _, course := range []string{"c", "c2"} {
		if err := put(s, course); err != nil {
			t.Fatal(err)
		}
	}
	back := load("c")
	if !reflect.DeepEqual(back, s) {
		t.Errorf("the story came back as\n%+v\nwant\n%+v", back, s)
	}
	out := filepath.Join(t.TempDir(), "story")
	if err := back.Write(out); err != nil {
		t.Fatal(err)
	}
	sameFiles(t, out, storyExample())
	for query, want := range map[string]int{
		`SELECT count(*) FROM story_page WHERE course_uuid = 'c'`:                                         2,
		`SELECT count(*) FROM story_page WHERE course_uuid = 'c' AND position = 1 AND eyebrow = 'APOGEE'`: 1,
		`SELECT count(*) FROM story_card WHERE course_uuid = 'c'`:                                         4,
		`SELECT count(*) FROM story_card WHERE course_uuid = 'c' AND subtitle IS NULL AND body IS NULL`:   2,
		`SELECT count(*) FROM chapter_story WHERE chapter_uuid = 'ch'`:                                    1,
		`SELECT count(*) FROM story_term WHERE course_uuid = 'c' AND linked = 0`:                          1,
	} {
		var n int
		if err := d.QueryRow(query).Scan(&n); err != nil || n != want {
			t.Errorf("%s = %d, want %d (%v)", query, n, want, err)
		}
	}
	if problems, err := db.Validate(context.Background(), d); err != nil || len(problems) > 0 {
		t.Error(problems, err)
	}

	small := Story{Pages: []Page{{Title: "Liftoff", Kind: "text", Action: "Watch it fly", Body: "Up it goes."}}}
	if err := put(small, "c"); err != nil {
		t.Fatal(err)
	}
	if got := load("c"); !reflect.DeepEqual(got, small) {
		t.Errorf("after a second Put: %+v, want %+v", got, small)
	}
	if got := load("c2"); !reflect.DeepEqual(got, s) {
		t.Error("a Put of one course's story changed another's")
	}
	if err := put(Story{}, "c"); err != nil {
		t.Fatal(err)
	}
	var rows int
	d.QueryRow(`SELECT (SELECT count(*) FROM story_page WHERE course_uuid = 'c') + (SELECT count(*) FROM story_card WHERE course_uuid = 'c') +
		(SELECT count(*) FROM chapter_story WHERE chapter_uuid = 'ch') + (SELECT count(*) FROM story_term WHERE course_uuid = 'c')`).Scan(&rows)
	if got := load("c"); !got.Empty() || rows != 0 {
		t.Errorf("after a Put of no story: %+v, and %d rows", got, rows)
	}
	err = put(Story{Chapters: []ChapterStory{{Chapter: "9-orbit", Briefing: "Up.", Debrief: "Around."}}}, "c")
	if err == nil || !strings.Contains(err.Error(), "the course has no chapter 9-orbit") {
		t.Errorf("Put of the story of a chapter the course does not have: %v", err)
	}
}

// A story another producer wrote, with gaps between its positions and its terms in any order,
// loads in the order the app shows it, and a package says it.
func TestLoadStoryOrders(t *testing.T) {
	d, err := db.Create(filepath.Join(t.TempDir(), "apogee.sqlite"), db.Info{Title: "Apogee", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, s := range []string{
		`INSERT INTO course VALUES ('c', 'learn-launches', 'Learn Launches', 'Course', 'Launch code.', 'An invented course.', 0)`,
		`INSERT INTO story_page VALUES ('c', 7, 'text', NULL, 'Later', 'The second page.', 'Next')`,
		`INSERT INTO story_page VALUES ('c', 2, 'text', NULL, 'First', 'The first page.', 'Next')`,
		`INSERT INTO story_card VALUES ('c', 7, 5, 'B', 'Second card', NULL, NULL), ('c', 7, 1, 'A', 'First card', NULL, NULL)`,
		`INSERT INTO story_term VALUES ('c', 'T-minus', 'The time left.', 1), ('c', 'pad', 'The concrete.', 0), ('c', 'AOS', 'Acquisition of signal.', 1)`,
	} {
		if _, err := d.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	s, err := LoadStory(d, "c")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, p := range s.Pages {
		got = append(got, p.Title)
		for _, c := range p.Cards {
			got = append(got, c.Title)
		}
	}
	for _, term := range s.Terms {
		got = append(got, term.Term)
	}
	if want := []string{"First", "Later", "First card", "Second card", "AOS", "pad", "T-minus"}; !slices.Equal(got, want) {
		t.Errorf("LoadStory gives %q, want %q", got, want)
	}
	out := filepath.Join(t.TempDir(), "story")
	if err := s.Write(out); err != nil {
		t.Fatal(err)
	}
	if again, err := ReadStory(out); err != nil || !reflect.DeepEqual(again, s) {
		t.Errorf("its package reads back as %+v, %v", again, err)
	}
}

const (
	onePage    = "# Welcome\n\nkind: text\naction: Go on\n\nHello.\n"
	oneChapter = "# 1-countdown\n\n## Briefing\n\nThe countdown starts.\n\n## Debrief\n\nIt ended well.\n"
	badgePage  = "# Your Badge\n\nkind: badge\neyebrow: APOGEE\naction: Print it\n\nWear it.\n\n## Flight Software Engineer\n\nlabel: Role\n"
)

// body is opening.md of one page whose body is text, on line 6 and on.
func body(text string) map[string]string {
	return map[string]string{"opening.md": "# Welcome\n\nkind: text\naction: Go on\n\n" + text + "\n"}
}

// withFields is opening.md of one page whose fields are lines, on line 3 and on.
func withFields(lines ...string) map[string]string {
	return map[string]string{"opening.md": "# Welcome\n\n" + strings.Join(lines, "\n") + "\n\nHello.\n"}
}

func TestReadStoryRefuses(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]string
		want  string
	}{
		// The folder and its files.
		{"a stray file", map[string]string{"opening.md": onePage, "notes.md": "to do\n"}, "notes.md does not belong in a story package, which holds opening.md, chapters.md and terms.md only"},
		{"text that is not UTF-8", map[string]string{"terms.md": "# pad\n\n\xff\n"}, "terms.md is not UTF-8 text"},
		{"an empty file", map[string]string{"opening.md": ""}, "opening.md is empty: leave it out when the course has none of it"},
		{"a file of blank lines", map[string]string{"terms.md": "\n\n"}, "terms.md is empty"},
		{"no newline at the end", map[string]string{"opening.md": strings.TrimSuffix(onePage, "\n")}, "opening.md line 6 has no newline at its end: the file ends right after its last line's newline"},
		{"a blank line at the end", map[string]string{"opening.md": onePage + "\n"}, "opening.md line 7 is blank: the file ends right after its last line's newline"},
		{"a blank line at the start", map[string]string{"chapters.md": "\n" + oneChapter}, "chapters.md line 1 is blank"},
		{"two blank lines", map[string]string{"opening.md": "# Welcome\n\n\nkind: text\naction: Go on\n\nHello.\n"}, "opening.md line 3 is a second blank line: exactly one comes between each two parts"},
		{"a line of spaces", map[string]string{"opening.md": "# Welcome\n  \nkind: text\naction: Go on\n\nHello.\n"}, "opening.md line 2 holds only spaces: a blank line is empty"},

		// opening.md: pages.
		{"no title line", map[string]string{"opening.md": strings.TrimPrefix(onePage, "# ")}, `opening.md line 1: "Welcome" is not a "# " line with the page's title`},
		{"no blank line after a title", map[string]string{"opening.md": "# Welcome\nkind: text\naction: Go on\n\nHello.\n"}, "opening.md line 2: a blank line comes between the page's title and the line after it"},
		{"an empty title", map[string]string{"opening.md": "# \n\nkind: text\naction: Go on\n\nHello.\n"}, "opening.md line 1: the page's title is empty"},
		{"a space before a title", map[string]string{"opening.md": "#  Welcome\n\nkind: text\naction: Go on\n\nHello.\n"}, `opening.md line 1: the page's title " Welcome" starts or ends with a space`},
		{"a picture in a title", map[string]string{"opening.md": "# ![logo](https://example.com/logo.png) Welcome\n\nkind: text\naction: Go on\n\nHello.\n"}, "opening.md line 1 holds 1 pictures, 0 videos, 0 links to videos and 0 widgets; the app shows none of them, so leave them out"},
		{"a page without fields", map[string]string{"opening.md": "# Welcome\n"}, `opening.md line 1: page "Welcome" has no fields; kind, eyebrow (optional) and action come after its title`},
		{"a body where the fields go", map[string]string{"opening.md": "# Welcome\n\nHello.\n"}, `opening.md line 3: "Hello." is not one of a page's fields, kind, eyebrow (optional) and action, which come right after its title`},
		{"no blank line after the fields", map[string]string{"opening.md": "# Welcome\n\nkind: text\naction: Go on\nHello.\n"}, `opening.md line 5: "Hello." is not one of a page's fields, kind, eyebrow (optional) and action; a blank line comes between the fields and what follows them`},
		{"a field a page does not have", withFields("kind: text", "colour: red", "action: Go on"), `opening.md line 4: "colour: red" is not one of a page's fields, kind, eyebrow (optional) and action`},
		{"fields out of order", withFields("kind: text", "action: Go on", "eyebrow: MONDAY"), "opening.md line 5: eyebrow comes after action; a page's fields are kind, eyebrow (optional) and action, in this order"},
		{"a field twice", withFields("kind: text", "kind: badge", "action: Go on"), "opening.md line 4: the page has kind twice"},
		{"no kind", withFields("eyebrow: MONDAY", "action: Go on"), "opening.md line 3: the page has no kind; its fields are kind, eyebrow (optional) and action, in this order"},
		{"no action", withFields("kind: text", "eyebrow: MONDAY"), "opening.md line 4: the page has no action"},
		{"an unknown kind", withFields("kind: video", "action: Go on"), `opening.md line 3: kind "video" is not text or badge`},
		{"an empty field", withFields("kind: text", "action:"), "opening.md line 4: action is empty"},
		{"no space after a colon", withFields("kind: text", "action:Go on"), `opening.md line 4: "action:Go on" has no space after its colon`},
		{"a space after a value", withFields("kind: text", "action: Go on "), `opening.md line 4: action "Go on " starts or ends with a space`},
		{"a page without a body", map[string]string{"opening.md": "# Welcome\n\nkind: text\naction: Go on\n"}, `opening.md line 1: page "Welcome" has no body; one or more paragraphs come after its fields`},
		{"cards but no body", map[string]string{"opening.md": "# Welcome\n\nkind: text\naction: Go on\n\n## Pad 2\n\nlabel: P2\n"}, `opening.md line 1: page "Welcome" has no body`},
		{"a field after the blank line", body("eyebrow: MONDAY\n\nHello."), `opening.md line 6: "eyebrow: MONDAY" starts like a field of the page, but the page's fields end at the blank line above it`},
		{"a badge page without cards", map[string]string{"opening.md": "# Your Badge\n\nkind: badge\naction: Print it\n\nWear it.\n"}, `opening.md line 1: badge page "Your Badge" has no cards; its cards are the badge's lines, one at least`},
		{"two badge pages", map[string]string{"opening.md": badgePage + "\n" + strings.Replace(badgePage, "Your Badge", "Another Badge", 1)}, `opening.md line 13: page "Another Badge" is a second badge page; an opening has at most one, and its first is at line 1`},

		// opening.md: cards.
		{"a card without fields", map[string]string{"opening.md": onePage + "\n## Pad 2\n"}, `opening.md line 8: card "Pad 2" has no fields; label and subtitle (optional) come after its title`},
		{"a card without a label", map[string]string{"opening.md": onePage + "\n## Pad 2\n\nsubtitle: The east pad\n"}, "opening.md line 10: the card has no label; its fields are label and subtitle (optional), in this order"},
		{"a card with an empty title", map[string]string{"opening.md": onePage + "\n## \n\nlabel: P2\n"}, "opening.md line 8: the card's title is empty"},
		{"a card's body of two paragraphs", map[string]string{"opening.md": onePage + "\n## Pad 2\n\nlabel: P2\n\nThe east pad.\n\nIt faces the sea.\n"}, `opening.md line 14: a second paragraph in card "Pad 2"; a card's body is one paragraph`},
		{"a card's field after the blank line", map[string]string{"opening.md": onePage + "\n## Pad 2\n\nlabel: P2\n\nsubtitle: The east pad\n"}, `opening.md line 12: "subtitle: The east pad" starts like a field of the card, but the card's fields end at the blank line above it`},

		// Paragraphs.
		{"a space before a line", body(" Hello."), `opening.md line 6: " Hello." starts or ends with a space`},
		{"a space after a line", body("Hello.\nAnd welcome. "), `opening.md line 7: "And welcome. " starts or ends with a space`},
		{"a list", body("- Hello."), `opening.md line 6: "- Hello." starts with "- ", as a list item does; a paragraph's lines start with text`},
		{"a list of stars", body("* Hello."), `opening.md line 6: "* Hello." starts with "* ", as a list item does`},
		{"a list of pluses", body("+ Hello."), `opening.md line 6: "+ Hello." starts with "+ ", as a list item does`},
		{"a numbered list", body("1. Hello."), `opening.md line 6: "1. Hello." starts with "1. ", as a numbered list item does`},
		{"a list numbered with parentheses", body("Hello.\n12) And welcome."), `opening.md line 7: "12) And welcome." starts with "12) ", as a numbered list item does`},
		{"a heading", body("Hello.\n# And welcome."), `opening.md line 7: "# And welcome." starts with "#", as a heading does`},
		{"a deeper heading", body("### Hello."), `opening.md line 6: "### Hello." starts with "#", as a heading does`},
		{"a quote", body("> Hello."), `opening.md line 6: "> Hello." starts with ">", as a quote does`},
		{"a code block", body("```go\nfmt.Println(\"Hello.\")\n```"), "opening.md line 6: \"```go\" starts with \"```\", as a code block does"},
		{"a code block of tildes", body("~~~\nHello.\n~~~"), `opening.md line 6: "~~~" starts with "~~~", as a code block does`},
		{"a table", body("| Pad | Rocket |"), `opening.md line 6: "| Pad | Rocket |" starts with "|", as a table row does`},
		{"HTML", body("<b>Hello.</b>"), `opening.md line 6: "<b>Hello.</b>" starts with "<", as HTML does`},
		{"HTML in a line", body("Press <kbd>Enter</kbd> now."), `opening.md line 6: "Press <kbd>Enter</kbd> now." holds HTML; a paragraph may use code spans, emphasis and strong emphasis, and no other Markdown`},
		{"a link", body("Read this.\nSee [the guide](https://example.com/guide) first."), `opening.md line 7: "See [the guide](https://example.com/guide) first." holds a link`},
		{"an autolink", body("See <https://example.com/guide> first."), `opening.md line 6: "See <https://example.com/guide> first." holds a link`},
		{"a line break", body("The first line\\\nand the second."), `opening.md line 6: "The first line\\" holds a backslash at its end, which breaks the line there`},
		{"a heading's underline", body("Hello.\n---"), `opening.md line 7: Markdown would not read "---" as part of a paragraph: a line like it makes a heading, a rule, a list item or a link definition`},
		{"a rule", body("Hello.\n***"), `opening.md line 7: Markdown would not read "***" as part of a paragraph`},
		{"an empty list item", body("Hello.\n\n-"), `opening.md line 8: Markdown would not read "-" as part of a paragraph`},
		{"a link definition", body("[guide]: https://example.com/guide"), `opening.md line 6: Markdown would not read "[guide]: https://example.com/guide" as part of a paragraph`},
		{"a picture", body("See ![the pad](https://example.com/pad.png) here."), "opening.md line 6 holds 1 pictures, 0 videos, 0 links to videos and 0 widgets; the app shows none of them, so leave them out"},
		{"a picture across two lines", body("See ![the\npad](https://example.com/pad.png) here."), "opening.md line 6 holds 1 pictures"},
		{"a picture in a card", map[string]string{"opening.md": onePage + "\n## Pad 2\n\nlabel: P2\n\nSee <img src=\"pad.png\" alt=\"the pad\"> here.\n"}, "opening.md line 12 holds 1 pictures"},
		{"a link to a video", body("Watch [the launch](https://youtu.be/abc123) now."), "opening.md line 6 holds 0 pictures, 0 videos, 1 links to videos and 0 widgets"},

		// chapters.md.
		{"a chapter's title that is not a slug", map[string]string{"chapters.md": strings.Replace(oneChapter, "1-countdown", "The Countdown", 1)}, `chapters.md line 1: "The Countdown" is not a chapter's slug: lowercase letters, digits and single hyphens`},
		{"a chapter twice", map[string]string{"chapters.md": oneChapter + "\n" + oneChapter}, "chapters.md line 11: chapter 1-countdown comes a second time; its story is at line 1"},
		{"a chapter without its story", map[string]string{"chapters.md": "# 1-countdown\n"}, `chapters.md line 1: chapter 1-countdown has no briefing; a "## Briefing" line and its paragraphs come next`},
		{"a briefing without its line", map[string]string{"chapters.md": "# 1-countdown\n\nThe countdown starts.\n"}, `chapters.md line 3: "The countdown starts." is not "## Briefing"; chapter 1-countdown's story is its briefing, then its debrief`},
		{"the debrief first", map[string]string{"chapters.md": "# 1-countdown\n\n## Debrief\n\nIt ended well.\n"}, `chapters.md line 3: "## Debrief" is not "## Briefing"`},
		{"no blank line after the briefing's line", map[string]string{"chapters.md": "# 1-countdown\n\n## Briefing\nThe countdown starts.\n\n## Debrief\n\nIt ended well.\n"}, `chapters.md line 4: a blank line comes after "## Briefing"`},
		{"an empty briefing", map[string]string{"chapters.md": "# 1-countdown\n\n## Briefing\n\n## Debrief\n\nIt ended well.\n"}, `chapters.md line 3: chapter 1-countdown's briefing is empty; one or more paragraphs come after "## Briefing"`},
		{"no debrief", map[string]string{"chapters.md": "# 1-countdown\n\n## Briefing\n\nThe countdown starts.\n"}, `chapters.md line 1: chapter 1-countdown has no debrief; a "## Debrief" line and its paragraphs come next`},
		{"the next chapter before the debrief", map[string]string{"chapters.md": "# 1-countdown\n\n## Briefing\n\nThe countdown starts.\n\n" + strings.Replace(oneChapter, "1-countdown", "2-orbit", 1)}, `chapters.md line 7: "# 2-orbit" is not "## Debrief"`},
		{"a part after the debrief", map[string]string{"chapters.md": oneChapter + "\n## Notes\n\nCheck the clock.\n"}, `chapters.md line 11: "## Notes" is not a "# " line with the chapter's slug`},
		{"a list in a debrief", map[string]string{"chapters.md": strings.Replace(oneChapter, "It ended well.", "- It ended well.", 1)}, `chapters.md line 9: "- It ended well." starts with "- ", as a list item does`},

		// terms.md.
		{"a term without its line", map[string]string{"terms.md": "pad\n\nThe concrete.\n"}, `terms.md line 1: "pad" is not a "# " line with the term`},
		{"terms out of order", map[string]string{"terms.md": "# pad\n\nThe concrete.\n\n# AOS\n\nAcquisition of signal.\n"}, `terms.md line 5: "AOS" comes after "pad" (line 1); the terms are in alphabetical order, ignoring case`},
		{"a term twice, ignoring case", map[string]string{"terms.md": "# AOS\n\nAcquisition of signal.\n\n# aos\n\nThe same again.\n"}, `terms.md line 5: "aos" is the term "AOS" of line 1 again, ignoring case; a term comes once`},
		{"linked: yes", map[string]string{"terms.md": "# pad\n\nlinked: yes\n\nThe concrete.\n"}, `terms.md line 3: "linked: yes" is not "linked: no"; a term is linked unless it says linked: no`},
		{"linked:no", map[string]string{"terms.md": "# pad\n\nlinked:no\n\nThe concrete.\n"}, `terms.md line 3: "linked:no" is not "linked: no"`},
		{"no blank line after linked: no", map[string]string{"terms.md": "# pad\n\nlinked: no\nThe concrete.\n"}, `terms.md line 4: a blank line comes after "linked: no"`},
		{"linked: no twice", map[string]string{"terms.md": "# pad\n\nlinked: no\n\nlinked: no\n\nThe concrete.\n"}, `terms.md line 5: "linked: no" starts like a field of the term, but the term's fields end at the blank line above it`},
		{"a term without a meaning", map[string]string{"terms.md": "# pad\n\n# T-minus\n\nThe time left.\n"}, `terms.md line 1: term "pad" has no meaning; one paragraph comes after it`},
		{"no meaning after linked: no", map[string]string{"terms.md": "# pad\n\nlinked: no\n"}, `terms.md line 1: term "pad" has no meaning`},
		{"a meaning of two paragraphs", map[string]string{"terms.md": "# pad\n\nThe concrete.\n\nA rocket stands on it.\n"}, `terms.md line 5: a second paragraph for term "pad"; a term's meaning is one paragraph`},
		{"a quote for a meaning", map[string]string{"terms.md": "# pad\n\n> The concrete.\n"}, `terms.md line 3: "> The concrete." starts with ">", as a quote does`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadStory(storyDir(t, tc.files))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("ReadStory: err = %v, want it to say %q", err, tc.want)
			}
		})
	}
}

// A folder, or a file that is not a regular one, has no place in a story package.
func TestReadStoryRefusesOtherEntries(t *testing.T) {
	folder := storyDir(t, nil)
	os.Mkdir(filepath.Join(folder, "drafts"), 0o755)
	misnamed := storyDir(t, nil)
	os.Mkdir(filepath.Join(misnamed, "terms.md"), 0o755)
	linked := storyDir(t, nil)
	elsewhere := filepath.Join(t.TempDir(), "opening.md")
	os.WriteFile(elsewhere, []byte(onePage), 0o644)
	if err := os.Symlink(elsewhere, filepath.Join(linked, "opening.md")); err != nil {
		t.Fatal(err)
	}
	for dir, want := range map[string]string{
		folder:                                "drafts does not belong in a story package",
		misnamed:                              "terms.md is not a regular file",
		linked:                                "opening.md is not a regular file",
		filepath.Join(t.TempDir(), "nowhere"): "no such file or directory",
	} {
		if _, err := ReadStory(dir); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("ReadStory(%s): err = %v, want it to say %q", dir, err, want)
		}
	}
}

// CheckChapters wants the chapters of the course only, in the course's order, and says where in
// chapters.md one is not.
func TestCheckChapters(t *testing.T) {
	s, err := parseStory(map[string]string{"chapters.md": oneChapter + "\n" + strings.Replace(oneChapter, "1-countdown", "3-landing", 1)})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		chapters []string
		want     string
	}{
		{[]string{"1-countdown", "2-orbit", "3-landing"}, ""},
		{[]string{"1-countdown", "2-orbit"}, "chapters.md line 11: course learn-launches has no chapter 3-landing; add it with author chapter, or leave its story out"},
		{[]string{"3-landing", "1-countdown"}, "chapters.md line 11: chapter 3-landing comes after 1-countdown, but the course has it before; chapters.md keeps the course's order"},
	} {
		err := s.CheckChapters("learn-launches", tc.chapters)
		if (err == nil) != (tc.want == "") || err != nil && err.Error() != tc.want {
			t.Errorf("CheckChapters(%q) = %v, want %q", tc.chapters, err, tc.want)
		}
	}
}

// Write refuses a story that its package cannot say, such as one another producer wrote
// into a library, and makes no folder for it.
func TestWriteStoryRefuses(t *testing.T) {
	badge := Page{Title: "Your Badge", Kind: "badge", Action: "Print it", Body: "Wear it.", Cards: []Card{{Title: "Flight Software Engineer", Label: "Role"}}}
	for _, tc := range []struct {
		name string
		s    Story
		want string
	}{
		{"a body that holds a card", Story{Pages: []Page{{Title: "Welcome", Kind: "text", Action: "Go on", Body: "Hello.\n\n## Pad 2\n\nlabel: P2"}}},
			"opening.md would read back as another story"},
		{"a title of two lines", Story{Pages: []Page{{Title: "Welcome\nto Apogee", Kind: "text", Action: "Go on", Body: "Hello."}}},
			"opening.md line 2: a blank line comes between the page's title and the line after it"},
		{"a badge page without cards", Story{Pages: []Page{{Title: "Your Badge", Kind: "badge", Action: "Print it", Body: "Wear it."}}},
			`opening.md line 1: badge page "Your Badge" has no cards`},
		{"two badge pages", Story{Pages: []Page{badge, badge}}, "is a second badge page"},
		{"a briefing that is a list", Story{Chapters: []ChapterStory{{Chapter: "1-countdown", Briefing: "- Count down.", Debrief: "Done."}}},
			`chapters.md line 5: "- Count down." starts with "- "`},
		{"terms alike ignoring case", Story{Terms: []Term{{Term: "AOS", Meaning: "Acquisition of signal.", Linked: true}, {Term: "aos", Meaning: "The same again.", Linked: true}}},
			`terms.md line 5: "aos" is the term "AOS" of line 1 again`},
		{"a meaning of two paragraphs", Story{Terms: []Term{{Term: "pad", Meaning: "The concrete.\n\nA rocket stands on it."}}},
			`a second paragraph for term "pad"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "story")
			err := tc.s.Write(dir)
			if err == nil || !strings.HasPrefix(err.Error(), "the story is not one a story package can say: ") || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Write: err = %v, want it to say %q", err, tc.want)
			}
			if _, err := os.Stat(dir); err == nil {
				t.Error("a refused Write made the folder")
			}
		})
	}
}

// The examples the help shows are what ReadStory takes, so writers can start from them.
func TestStoryLayoutExamplesAreValid(t *testing.T) {
	opening := "# Welcome to Apogee\n\nkind: text\neyebrow: PORT ELLERY · MONDAY 1 MARCH 2027\naction: Start my first day\n\n" +
		"Right now, eight satellites circle the Earth.\n\nToday you join the team that keeps them talking to the ground.\n\n" +
		"## Nkem Halloran\n\nlabel: NH\nsubtitle: Staff engineer, your onboarding buddy\n\nReviews your code: short, kind and exact.\n\n" +
		"# First, Your Badge\n\nkind: badge\neyebrow: APOGEE\naction: Print my badge\n\nEveryone at Apogee wears one.\n\n" +
		"## Ground Software Engineer\n\nlabel: Role\n"
	chapters := "# first-contact\n\n## Briefing\n\nDay one. By tradition, your first program writes the message we send up to Pipit-1.\n\n" +
		"## Debrief\n\nYour words went up and came back from space.\n"
	terms := "# AOS\n\nAcquisition of signal: the moment a satellite rises into a dish's view.\n\n" +
		"# pass\n\nlinked: no\n\nOne flyover of one satellite above one dish.\n"
	for name, example := range map[string]string{"opening.md": opening, "chapters.md": chapters, "terms.md": terms} {
		if !strings.Contains(StoryLayout, indent(example)) {
			t.Errorf("StoryLayout does not show the %s example", name)
		}
	}
	s, err := parseStory(map[string]string{"opening.md": opening, "chapters.md": chapters, "terms.md": terms})
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Summary(); got != "2 pages, 2 cards, 1 chapter and 2 terms" || s.Pages[1].Kind != "badge" || s.Pages[1].Eyebrow != "APOGEE" ||
		s.Pages[0].Cards[0].Body != "Reviews your code: short, kind and exact." || !s.Terms[0].Linked || s.Terms[1].Linked {
		t.Errorf("the examples read as %s: %+v", got, s)
	}
}

// Whatever a story file holds, reading it either refuses it or gives a story that is written
// back as the same text: a package has exactly one way to be written.
func FuzzStoryFiles(f *testing.F) {
	for _, name := range storyFiles {
		b, err := os.ReadFile(filepath.Join(storyExample(), name))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(b))
	}
	f.Add(onePage + "\n" + badgePage)
	f.Fuzz(func(t *testing.T, text string) {
		if pages, err := parseOpening(text); err == nil {
			if back := openingText(pages); back != text {
				t.Errorf("opening.md %q is written back as %q", text, back)
			}
			if _, err := (Story{Pages: pages}).files(); err != nil {
				t.Errorf("opening.md %q reads, but its story does not write: %v", text, err)
			}
		}
		if chapters, err := parseChapters(text); err == nil {
			if back := chaptersText(chapters); back != text {
				t.Errorf("chapters.md %q is written back as %q", text, back)
			}
			if _, err := (Story{Chapters: chapters}).files(); err != nil {
				t.Errorf("chapters.md %q reads, but its story does not write: %v", text, err)
			}
		}
		if terms, err := parseTerms(text); err == nil {
			if back := termsText(terms); back != text {
				t.Errorf("terms.md %q is written back as %q", text, back)
			}
			if _, err := (Story{Terms: terms}).files(); err != nil {
				t.Errorf("terms.md %q reads, but its story does not write: %v", text, err)
			}
		}
	})
}
