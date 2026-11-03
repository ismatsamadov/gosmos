package authoring

// A story package is the story a course is told in, as writers author it: opening.md, the
// pages the app shows before the course's first lesson; chapters.md, the chapters' briefings
// and debriefs; and terms.md, the field guide. StoryLayout is its specification. ReadStory
// reads one and Story.Write writes one; LoadStory and Story.Put map a story to and from a
// library's rows (story_page, story_card, chapter_story and story_term).
//
// Each file is a series of chunks, runs of lines with exactly one blank line between each two,
// and the parsers below read them in order. Every refusal names the file and the line.

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	mdtext "github.com/yuin/goldmark/text"

	"github.com/jobloom/gosmos/internal/store"
)

// StoryLayout is the story package specification, as the author tool's help shows it.
const StoryLayout = `A story package is a folder that holds a course's story, in up to three files:

  opening.md   the opening: the pages the app shows before the course's first lesson
  chapters.md  the chapters' briefings, which open them, and debriefs, which the app shows
               once a chapter's required lessons are done
  terms.md     the field guide: the words of the course's world, each with its meaning in
               plain language

A file left out means the course has none of that part, and a folder with none of them
clears the course's story; leave a file out rather than leave it empty. Nothing else belongs
in the folder. Files and folders whose names start with a dot are ignored. Every file is
UTF-8 text, and none is a symbolic link.

The examples below are indented by two spaces to set them apart. In the files, every line
starts at the left margin.

Each file is a series of parts with exactly one blank line between each two, and it ends
right after its last line's newline. A title, a slug, a term and a field's value are one
line each, never empty, and none starts or ends with a space.

A paragraph is one or more lines of text. It may use code spans, emphasis and strong
emphasis, and no other Markdown: no headings, lists, quotes, code blocks, tables, HTML,
links, line breaks (a backslash at a line's end) or pictures. No line starts or ends with a
space, and none starts with #, >, "- ", "* ", "+ ", a number followed by ". " or ") ", ` + "```" + `,
~~~, | or <. Each text is kept as written: a paragraph's lines joined by newlines, and the
paragraphs of a text by a blank line.

opening.md holds one or more pages, which the app shows in order:

  # Welcome to Apogee

  kind: text
  eyebrow: PORT ELLERY · MONDAY 1 MARCH 2027
  action: Start my first day

  Right now, eight satellites circle the Earth.

  Today you join the team that keeps them talking to the ground.

  ## Nkem Halloran

  label: NH
  subtitle: Staff engineer, your onboarding buddy

  Reviews your code: short, kind and exact.

  # First, Your Badge

  kind: badge
  eyebrow: APOGEE
  action: Print my badge

  Everyone at Apogee wears one.

  ## Ground Software Engineer

  label: Role

A page is a "# " line with its title; its fields; its body, one or more paragraphs; and
then its cards, if it has any. Its fields come one per line, each as its name, a colon, a
space and its value, in this order:

  kind     text, or badge for the page that shows the learner's badge
  eyebrow  optional: a short line above the title, such as a place and a day; on a badge
           page it also heads the badge, naming the organization the learner works for
  action   the button that moves on; the last page's opens the course

A card is a "## " line with its title; its fields, label and then subtitle if it has one;
and then, if it has one, its body: one paragraph.

  label     a short mark: initials, a station code, a line's name
  subtitle  optional: a line below the title

A text page shows its cards, the people and places it introduces, below its body. A badge
page shows the learner's badge: its eyebrow heads the badge (a badge page without one shows
no header), each of its cards is a line of the badge, and a field takes the name the team
calls the learner by. A badge page has at least one card, and an opening has at most one
badge page.

chapters.md holds the briefing and the debrief of any of the course's chapters, each
chapter at most once and in the course's order:

  # first-contact

  ## Briefing

  Day one. By tradition, your first program writes the message we send up to Pipit-1.

  ## Debrief

  Your words went up and came back from space.

A chapter is a "# " line with its slug; a "## Briefing" line, then one or more paragraphs,
which open the chapter; and a "## Debrief" line, then one or more paragraphs, which the app
shows once the chapter's required lessons are done.

terms.md holds the field guide's terms in alphabetical order, ignoring case, no two alike
ignoring case:

  # AOS

  Acquisition of signal: the moment a satellite rises into a dish's view.

  # pass

  linked: no

  One flyover of one satellite above one dish.

A term is a "# " line with the term, as readmes write it, then its meaning: one paragraph.
The first time a readme uses a term, outside code, the app links it to its meaning. A term
that lessons also use as an everyday word, such as pass, is never linked: it says
"linked: no", on a line of its own between the term and its meaning.

A body or a meaning never starts with one of its part's field names and a colon (kind:,
eyebrow: or action: on a page, label: or subtitle: on a card, linked: on a term): a field
comes with the others, before the blank line that ends them.`

// storyFiles are the files a story package may hold, in the order the help names them.
var storyFiles = []string{"opening.md", "chapters.md", "terms.md"}

// Story is the story a course is told in, as its story package holds it. A part the course
// does not have is empty.
type Story struct {
	Pages    []Page         // opening.md: the opening, in the order the app shows it
	Chapters []ChapterStory // chapters.md: in the course's order
	Terms    []Term         // terms.md: in alphabetical order, ignoring case
}

// Page is a page of the opening.
type Page struct {
	Title   string
	Kind    string // text or badge
	Eyebrow string // a short line above the title, or "" for none; a badge page's heads the badge
	Action  string // the button that moves on
	Body    string // one or more paragraphs, a blank line between each two
	Cards   []Card
}

// Card is a card of a page: a person or place a text page introduces, or a line of a badge.
type Card struct {
	Title    string
	Label    string // a short mark: initials, a station code, a line's name
	Subtitle string // or "" for none
	Body     string // one paragraph, or "" for none
}

// ChapterStory is a chapter's part in the story.
type ChapterStory struct {
	Chapter  string // the chapter's slug
	Briefing string // one or more paragraphs, a blank line between each two
	Debrief  string // one or more paragraphs, a blank line between each two
}

// Term is a word of the course's world, with its meaning.
type Term struct {
	Term    string // as readmes write it
	Meaning string // one paragraph
	Linked  bool   // readmes link to the meaning; false for a term lessons also use as an everyday word
}

// Empty reports whether the story has none of its parts.
func (s Story) Empty() bool { return len(s.Pages) == 0 && len(s.Chapters) == 0 && len(s.Terms) == 0 }

// Summary says what the story holds, such as "2 pages, 3 cards, 1 chapter and 4 terms", or
// "no story".
func (s Story) Summary() string {
	cards := 0
	for _, p := range s.Pages {
		cards += len(p.Cards)
	}
	var parts []string
	for _, c := range []struct {
		n     int
		thing string
	}{{len(s.Pages), "page"}, {cards, "card"}, {len(s.Chapters), "chapter"}, {len(s.Terms), "term"}} {
		if c.n > 0 {
			parts = append(parts, count(c.n, c.thing))
		}
	}
	switch len(parts) {
	case 0:
		return "no story"
	case 1:
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

// Changed names the files of the story package in which two stories differ, in the order the
// help names them.
func (s Story) Changed(t Story) []string {
	var out []string
	if !slices.EqualFunc(s.Pages, t.Pages, samePage) {
		out = append(out, "opening.md")
	}
	if !slices.Equal(s.Chapters, t.Chapters) {
		out = append(out, "chapters.md")
	}
	if !slices.Equal(s.Terms, t.Terms) {
		out = append(out, "terms.md")
	}
	return out
}

// samePage reports whether two pages are the same, their cards included.
func samePage(p, q Page) bool {
	return p.Title == q.Title && p.Kind == q.Kind && p.Eyebrow == q.Eyebrow && p.Action == q.Action &&
		p.Body == q.Body && slices.Equal(p.Cards, q.Cards)
}

// ReadStory reads and checks the story package in dir. It does not know the course's
// chapters: CheckChapters checks chapters.md against them.
func ReadStory(dir string) (Story, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Story{}, err
	}
	files := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		switch {
		case strings.HasPrefix(name, "."):
			continue
		case !slices.Contains(storyFiles, name):
			return Story{}, fmt.Errorf("%s does not belong in a story package, which holds opening.md, chapters.md and terms.md only", name)
		case !e.Type().IsRegular():
			return Story{}, fmt.Errorf("%s is not a regular file", name)
		}
		if files[name], err = readText(filepath.Join(dir, name)); err != nil {
			return Story{}, err
		}
	}
	return parseStory(files)
}

// parseStory reads a story from the texts of its package's files, by name.
func parseStory(files map[string]string) (Story, error) {
	var s Story
	var err error
	if text, ok := files["opening.md"]; ok {
		if s.Pages, err = parseOpening(text); err != nil {
			return Story{}, err
		}
	}
	if text, ok := files["chapters.md"]; ok {
		if s.Chapters, err = parseChapters(text); err != nil {
			return Story{}, err
		}
	}
	if text, ok := files["terms.md"]; ok {
		if s.Terms, err = parseTerms(text); err != nil {
			return Story{}, err
		}
	}
	return s, nil
}

// CheckChapters checks chapters.md against the course's chapters, given as their slugs in the
// course's order: the story tells only of chapters the course has, in the same order. course
// is the course's slug, for errors.
func (s Story) CheckChapters(course string, chapters []string) error {
	lines := headingLines(chaptersText(s.Chapters))
	prev := -1
	for i, c := range s.Chapters {
		k := slices.Index(chapters, c.Chapter)
		switch {
		case k < 0:
			return fmt.Errorf("chapters.md line %d: course %s has no chapter %s; add it with author chapter, or leave its story out", lines[i], course, c.Chapter)
		case k < prev:
			return fmt.Errorf("chapters.md line %d: chapter %s comes after %s, but the course has it before; chapters.md keeps the course's order", lines[i], c.Chapter, s.Chapters[i-1].Chapter)
		}
		prev = k
	}
	return nil
}

// headingLines returns the numbers of the "# " lines of a story file, in order.
func headingLines(text string) []int {
	var out []int
	for i, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "# ") {
			out = append(out, i+1)
		}
	}
	return out
}

// ---- Chunks ----

// chunk is a run of lines with no blank line among them: a title, a part's fields or a
// paragraph.
type chunk struct {
	line  int // the number of its first line, counted from 1
	lines []string
}

// chunks splits the text of a story file into its chunks. Exactly one blank line comes between
// each two, and the file ends right after its last line's newline.
func chunks(file, text string) ([]chunk, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("%s is empty: leave it out when the course has none of it", file)
	}
	lines := strings.Split(text, "\n")
	if lines[len(lines)-1] != "" {
		return nil, fmt.Errorf("%s line %d has no newline at its end: the file ends right after its last line's newline", file, len(lines))
	}
	lines = lines[:len(lines)-1]
	var out []chunk
	for i, line := range lines {
		n := i + 1
		switch {
		case line == "" && i == len(lines)-1:
			return nil, fmt.Errorf("%s line %d is blank: the file ends right after its last line's newline", file, n)
		case line == "" && i == 0:
			return nil, fmt.Errorf("%s line 1 is blank: the file starts with its first part", file)
		case line == "" && lines[i-1] == "":
			return nil, fmt.Errorf("%s line %d is a second blank line: exactly one comes between each two parts", file, n)
		case line == "":
		case strings.TrimSpace(line) == "":
			return nil, fmt.Errorf("%s line %d holds only spaces: a blank line is empty", file, n)
		case i > 0 && lines[i-1] != "":
			out[len(out)-1].lines = append(out[len(out)-1].lines, line)
		default:
			out = append(out, chunk{line: n, lines: []string{line}})
		}
	}
	return out, nil
}

// isHeading reports whether a chunk starts with a "# " or "## " line, a part's title.
func isHeading(c chunk) bool {
	return strings.HasPrefix(c.lines[0], "# ") || strings.HasPrefix(c.lines[0], "## ")
}

// heading reads a chunk that is one line: marker ("# " or "## ") and the text it marks, which
// what names in errors, such as "the page's title".
func heading(file string, c chunk, marker, what string) (string, error) {
	t, ok := strings.CutPrefix(c.lines[0], marker)
	if !ok {
		return "", fmt.Errorf("%s line %d: %q is not a %q line with %s", file, c.line, c.lines[0], marker, what)
	}
	if len(c.lines) > 1 {
		return "", fmt.Errorf("%s line %d: a blank line comes between %s and the line after it", file, c.line+1, what)
	}
	return t, checkValue(file, c.line, what, t)
}

// checkValue checks a one-line value of a story file, which what names in errors: it is not
// empty, it has no space at either end, and it holds no pictures.
func checkValue(file string, line int, what, v string) error {
	switch {
	case v == "":
		return fmt.Errorf("%s line %d: %s is empty", file, line, what)
	case strings.TrimSpace(v) != v:
		return fmt.Errorf("%s line %d: %s %q starts or ends with a space", file, line, what, v)
	}
	return mediaError(fmt.Sprintf("%s line %d", file, line), v)
}

// field is one of a part's fields, a line of its name, a colon, a space and its value.
type field struct {
	name     string
	optional bool
}

var (
	pageFields = []field{{"kind", false}, {"eyebrow", true}, {"action", false}}
	cardFields = []field{{"label", false}, {"subtitle", true}}
	termFields = []field{{"linked", true}}
)

// reFieldName matches what reads as a field's name, such as colour in "colour: red", rather
// than text, such as "Right now" in "Right now: launch".
var reFieldName = regexp.MustCompile(`^[a-z]+$`)

// fields reads a chunk of a part's fields, which come in the order want gives them. part names
// the part in errors, such as "page".
func fields(file string, c chunk, part string, want []field) (map[string]string, error) {
	got := map[string]string{}
	next := 0 // want[next:] are the fields that may still come
	for i, line := range c.lines {
		n := c.line + i
		name, rest, ok := strings.Cut(line, ":")
		k := slices.IndexFunc(want, func(f field) bool { return f.name == name })
		if !ok || k < 0 {
			err := fmt.Sprintf("%s line %d: %q is not one of a %s's fields, %s", file, n, line, part, fieldList(want))
			switch {
			case ok && reFieldName.MatchString(name): // a field the part does not have
			case i > 0:
				err += "; a blank line comes between the fields and what follows them"
			default:
				err += ", which come right after its title"
			}
			return nil, errors.New(err)
		}
		if _, dup := got[name]; dup {
			return nil, fmt.Errorf("%s line %d: the %s has %s twice", file, n, part, name)
		}
		if k < next {
			return nil, fmt.Errorf("%s line %d: %s comes after %s; a %s's fields are %s, in this order", file, n, name, want[next-1].name, part, fieldList(want))
		}
		if missing := firstRequired(want[next:k]); missing != "" {
			return nil, fmt.Errorf("%s line %d: the %s has no %s; its fields are %s, in this order", file, n, part, missing, fieldList(want))
		}
		v, spaced := strings.CutPrefix(rest, " ")
		if !spaced && rest != "" {
			return nil, fmt.Errorf("%s line %d: %q has no space after its colon", file, n, line)
		}
		if err := checkValue(file, n, name, v); err != nil {
			return nil, err
		}
		got[name] = v
		next = k + 1
	}
	if missing := firstRequired(want[next:]); missing != "" {
		return nil, fmt.Errorf("%s line %d: the %s has no %s; its fields are %s, in this order", file, c.line+len(c.lines)-1, part, missing, fieldList(want))
	}
	return got, nil
}

// firstRequired returns the name of the first field of fs that is not optional, or "".
func firstRequired(fs []field) string {
	for _, f := range fs {
		if !f.optional {
			return f.name
		}
	}
	return ""
}

// fieldList names a part's fields in their order, such as "label and subtitle (optional)".
func fieldList(want []field) string {
	var names []string
	for _, f := range want {
		if f.optional {
			names = append(names, f.name+" (optional)")
		} else {
			names = append(names, f.name)
		}
	}
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// cutOff refuses a part's body or meaning whose first line starts like one of the part's
// fields: a field that a blank line cut off from the others.
func cutOff(file string, c chunk, part string, want []field) error {
	for _, f := range want {
		if strings.HasPrefix(c.lines[0], f.name+":") {
			return fmt.Errorf("%s line %d: %q starts like a field of the %s, but the %s's fields end at the blank line above it", file, c.line, c.lines[0], part, part)
		}
	}
	return nil
}

// ---- Paragraphs ----

var reListNumber = regexp.MustCompile(`^[0-9]+[.)] `)

// blockMarks are how a line starts that Markdown makes something other than a paragraph of,
// and what it makes.
var blockMarks = []struct{ mark, what string }{
	{"#", "a heading"}, {">", "a quote"}, {"- ", "a list item"}, {"* ", "a list item"}, {"+ ", "a list item"},
	{"```", "a code block"}, {"~~~", "a code block"}, {"|", "a table row"}, {"<", "HTML"},
}

// blockMark returns how a line starts that a paragraph's lines never start with, and what
// Markdown makes of it, or "" for a line a paragraph may hold.
func blockMark(line string) (mark, what string) {
	if m := reListNumber.FindString(line); m != "" {
		return m, "a numbered list item"
	}
	for _, b := range blockMarks {
		if strings.HasPrefix(line, b.mark) {
			return b.mark, b.what
		}
	}
	return "", ""
}

// paragraph reads a chunk that is one paragraph of a story file and returns its text, its lines
// joined by newlines.
func paragraph(file string, c chunk) (string, error) {
	for i, line := range c.lines {
		n := c.line + i
		if strings.TrimSpace(line) != line {
			return "", fmt.Errorf("%s line %d: %q starts or ends with a space", file, n, line)
		}
		if mark, what := blockMark(line); mark != "" {
			return "", fmt.Errorf("%s line %d: %q starts with %q, as %s does; a paragraph's lines start with text", file, n, line, mark, what)
		}
		if err := mediaError(fmt.Sprintf("%s line %d", file, n), line); err != nil {
			return "", err
		}
	}
	text := strings.Join(c.lines, "\n")
	if err := mediaError(fmt.Sprintf("%s line %d", file, c.line), text); err != nil {
		return "", err // a picture written across lines
	}
	return text, inlineOnly(file, c)
}

// markdown reads a paragraph as CommonMark, without extensions, as media reads readmes.
var markdown = goldmark.New().Parser()

// inlineOnly checks that Markdown reads a chunk's lines as one paragraph, made of text, code
// spans, emphasis and strong emphasis only.
func inlineOnly(file string, c chunk) error {
	// The first line that keeps the lines up to it from being one paragraph: a line of
	// dashes that makes a heading of the line above, a rule, an empty list item, a link
	// definition.
	for k := range c.lines {
		if !oneParagraph(c.lines[:k+1]) {
			return fmt.Errorf("%s line %d: Markdown would not read %q as part of a paragraph: a line like it makes a heading, a rule, a list item or a link definition",
				file, c.line+k, c.lines[k])
		}
	}
	src := []byte(strings.Join(c.lines, "\n"))
	var bad error
	_ = ast.Walk(markdown.Parse(mdtext.NewReader(src)).FirstChild(), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		var what string
		switch n := n.(type) {
		case *ast.Paragraph, *ast.Emphasis, *ast.String:
			return ast.WalkContinue, nil
		case *ast.CodeSpan:
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			if !n.HardLineBreak() {
				return ast.WalkContinue, nil
			}
			what = "a backslash at its end, which breaks the line there"
		case *ast.Link, *ast.AutoLink:
			what = "a link"
		case *ast.Image:
			what = "a picture"
		case *ast.RawHTML:
			what = "HTML"
		default:
			what = "Markdown of the kind " + n.Kind().String()
		}
		k := bytes.Count(src[:max(n.Pos(), 0)], []byte("\n"))
		bad = fmt.Errorf("%s line %d: %q holds %s; a paragraph may use code spans, emphasis and strong emphasis, and no other Markdown",
			file, c.line+k, c.lines[k], what)
		return ast.WalkStop, nil
	})
	return bad
}

// oneParagraph reports whether Markdown reads lines as one paragraph that holds them all.
func oneParagraph(lines []string) bool {
	doc := markdown.Parse(mdtext.NewReader([]byte(strings.Join(lines, "\n"))))
	p, ok := doc.FirstChild().(*ast.Paragraph)
	return ok && p.NextSibling() == nil && p.Lines().Len() == len(lines)
}

// paragraphsFrom reads the paragraphs from chunk i up to the next "# " or "## " line, and
// returns them and the chunk after them.
func paragraphsFrom(file string, cs []chunk, i int) ([]string, int, error) {
	var out []string
	for ; i < len(cs) && !isHeading(cs[i]); i++ {
		p, err := paragraph(file, cs[i])
		if err != nil {
			return nil, i, err
		}
		out = append(out, p)
	}
	return out, i, nil
}

// ---- opening.md ----

// parseOpening reads opening.md: one or more pages, each with its cards.
func parseOpening(text string) ([]Page, error) {
	const file = "opening.md"
	cs, err := chunks(file, text)
	if err != nil {
		return nil, err
	}
	var pages []Page
	badge := 0 // the line of the badge page's title, once there is one
	for i := 0; i < len(cs); {
		at := cs[i].line
		title, err := heading(file, cs[i], "# ", "the page's title")
		if err != nil {
			return nil, err
		}
		p := Page{Title: title}
		if i++; i == len(cs) {
			return nil, fmt.Errorf("%s line %d: page %q has no fields; %s come after its title", file, at, title, fieldList(pageFields))
		}
		f, err := fields(file, cs[i], "page", pageFields)
		if err != nil {
			return nil, err
		}
		p.Kind, p.Eyebrow, p.Action = f["kind"], f["eyebrow"], f["action"]
		if p.Kind != "text" && p.Kind != "badge" {
			return nil, fmt.Errorf("%s line %d: kind %q is not text or badge", file, cs[i].line, p.Kind)
		}
		if i++; i < len(cs) && !isHeading(cs[i]) {
			if err := cutOff(file, cs[i], "page", pageFields); err != nil {
				return nil, err
			}
		}
		body, next, err := paragraphsFrom(file, cs, i)
		if err != nil {
			return nil, err
		}
		if len(body) == 0 {
			return nil, fmt.Errorf("%s line %d: page %q has no body; one or more paragraphs come after its fields", file, at, title)
		}
		p.Body, i = strings.Join(body, "\n\n"), next
		for i < len(cs) && strings.HasPrefix(cs[i].lines[0], "## ") {
			var c Card
			if c, i, err = parseCard(file, cs, i); err != nil {
				return nil, err
			}
			p.Cards = append(p.Cards, c)
		}
		if p.Kind == "badge" {
			if len(p.Cards) == 0 {
				return nil, fmt.Errorf("%s line %d: badge page %q has no cards; its cards are the badge's lines, one at least", file, at, title)
			}
			if badge > 0 {
				return nil, fmt.Errorf("%s line %d: page %q is a second badge page; an opening has at most one, and its first is at line %d", file, at, title, badge)
			}
			badge = at
		}
		pages = append(pages, p)
	}
	return pages, nil
}

// parseCard reads the card whose "## " line is chunk i, and returns it and the chunk after it.
func parseCard(file string, cs []chunk, i int) (Card, int, error) {
	at := cs[i].line
	title, err := heading(file, cs[i], "## ", "the card's title")
	if err != nil {
		return Card{}, i, err
	}
	c := Card{Title: title}
	if i++; i == len(cs) {
		return Card{}, i, fmt.Errorf("%s line %d: card %q has no fields; %s come after its title", file, at, title, fieldList(cardFields))
	}
	f, err := fields(file, cs[i], "card", cardFields)
	if err != nil {
		return Card{}, i, err
	}
	c.Label, c.Subtitle = f["label"], f["subtitle"]
	if i++; i == len(cs) || isHeading(cs[i]) {
		return c, i, nil
	}
	if err := cutOff(file, cs[i], "card", cardFields); err != nil {
		return Card{}, i, err
	}
	if c.Body, err = paragraph(file, cs[i]); err != nil {
		return Card{}, i, err
	}
	if i++; i < len(cs) && !isHeading(cs[i]) {
		return Card{}, i, fmt.Errorf("%s line %d: a second paragraph in card %q; a card's body is one paragraph", file, cs[i].line, title)
	}
	return c, i, nil
}

// ---- chapters.md ----

// parseChapters reads chapters.md: the briefing and the debrief of chapters, each at most once.
func parseChapters(text string) ([]ChapterStory, error) {
	const file = "chapters.md"
	cs, err := chunks(file, text)
	if err != nil {
		return nil, err
	}
	var out []ChapterStory
	seen := map[string]int{} // the line of each chapter's slug
	for i := 0; i < len(cs); {
		at := cs[i].line
		slug, err := heading(file, cs[i], "# ", "the chapter's slug")
		if err != nil {
			return nil, err
		}
		if !ValidSlug(slug) {
			return nil, fmt.Errorf("%s line %d: %q is not a chapter's slug: lowercase letters, digits and single hyphens", file, at, slug)
		}
		if first, dup := seen[slug]; dup {
			return nil, fmt.Errorf("%s line %d: chapter %s comes a second time; its story is at line %d", file, at, slug, first)
		}
		seen[slug] = at
		c := ChapterStory{Chapter: slug}
		if c.Briefing, i, err = section(file, cs, i+1, at, slug, "Briefing"); err != nil {
			return nil, err
		}
		if c.Debrief, i, err = section(file, cs, i, at, slug, "Debrief"); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

// section reads a chapter's briefing or debrief (name), starting at chunk i: a "## " line with
// its name, then one or more paragraphs. at is the line of the chapter's slug. It returns the
// paragraphs and the chunk after them.
func section(file string, cs []chunk, i, at int, chapter, name string) (string, int, error) {
	head, part := "## "+name, strings.ToLower(name)
	switch {
	case i == len(cs):
		return "", i, fmt.Errorf("%s line %d: chapter %s has no %s; a %q line and its paragraphs come next", file, at, chapter, part, head)
	case cs[i].lines[0] != head:
		return "", i, fmt.Errorf("%s line %d: %q is not %q; chapter %s's story is its briefing, then its debrief", file, cs[i].line, cs[i].lines[0], head, chapter)
	case len(cs[i].lines) > 1:
		return "", i, fmt.Errorf("%s line %d: a blank line comes after %q", file, cs[i].line+1, head)
	}
	paras, next, err := paragraphsFrom(file, cs, i+1)
	if err != nil {
		return "", next, err
	}
	if len(paras) == 0 {
		return "", next, fmt.Errorf("%s line %d: chapter %s's %s is empty; one or more paragraphs come after %q", file, cs[i].line, chapter, part, head)
	}
	return strings.Join(paras, "\n\n"), next, nil
}

// ---- terms.md ----

// parseTerms reads terms.md: the field guide's terms, in alphabetical order ignoring case, each
// with its meaning.
func parseTerms(text string) ([]Term, error) {
	const file = "terms.md"
	cs, err := chunks(file, text)
	if err != nil {
		return nil, err
	}
	var out []Term
	last := 0 // the line of the term before
	for i := 0; i < len(cs); {
		at := cs[i].line
		term, err := heading(file, cs[i], "# ", "the term")
		if err != nil {
			return nil, err
		}
		if len(out) > 0 {
			prev := out[len(out)-1].Term
			switch c := compareTerms(prev, term); {
			case c == 0:
				return nil, fmt.Errorf("%s line %d: %q is the term %q of line %d again, ignoring case; a term comes once", file, at, term, prev, last)
			case c > 0:
				return nil, fmt.Errorf("%s line %d: %q comes after %q (line %d); the terms are in alphabetical order, ignoring case", file, at, term, prev, last)
			}
		}
		t := Term{Term: term, Linked: true}
		if i++; i < len(cs) && strings.HasPrefix(cs[i].lines[0], "linked:") {
			switch {
			case cs[i].lines[0] != "linked: no":
				return nil, fmt.Errorf("%s line %d: %q is not \"linked: no\"; a term is linked unless it says linked: no", file, cs[i].line, cs[i].lines[0])
			case len(cs[i].lines) > 1:
				return nil, fmt.Errorf("%s line %d: a blank line comes after \"linked: no\"", file, cs[i].line+1)
			}
			t.Linked = false
			i++
		}
		if i == len(cs) || isHeading(cs[i]) {
			return nil, fmt.Errorf("%s line %d: term %q has no meaning; one paragraph comes after it", file, at, term)
		}
		if err := cutOff(file, cs[i], "term", termFields); err != nil {
			return nil, err
		}
		if t.Meaning, err = paragraph(file, cs[i]); err != nil {
			return nil, err
		}
		if i++; i < len(cs) && !isHeading(cs[i]) {
			return nil, fmt.Errorf("%s line %d: a second paragraph for term %q; a term's meaning is one paragraph", file, cs[i].line, term)
		}
		out = append(out, t)
		last = at
	}
	return out, nil
}

// compareTerms orders terms as terms.md does: alphabetically, ignoring case.
func compareTerms(a, b string) int { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) }

// ---- Writing ----

// Write writes the package into dir, which must not exist yet or be empty; a story without
// parts is an empty folder. It refuses a story that a package cannot say exactly, so that
// reading what it wrote gives the story back.
func (s Story) Write(dir string) error {
	if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
		return fmt.Errorf("%s is not empty", dir)
	}
	files, err := s.files()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, name := range sortedKeys(files) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(files[name]), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// files returns the texts of the story's package, by file name, once it has read them back as
// the same story.
func (s Story) files() (map[string]string, error) {
	files := map[string]string{}
	if len(s.Pages) > 0 {
		files["opening.md"] = openingText(s.Pages)
	}
	if len(s.Chapters) > 0 {
		files["chapters.md"] = chaptersText(s.Chapters)
	}
	if len(s.Terms) > 0 {
		files["terms.md"] = termsText(s.Terms)
	}
	back, err := parseStory(files)
	if changed := back.Changed(s); err == nil && len(changed) > 0 {
		err = fmt.Errorf("%s would read back as another story", strings.Join(changed, ", "))
	}
	if err != nil {
		return nil, fmt.Errorf("the story is not one a story package can say: %w", err)
	}
	return files, nil
}

func openingText(pages []Page) string {
	var b strings.Builder
	for i, p := range pages {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("# " + p.Title + "\n\nkind: " + p.Kind + "\n")
		if p.Eyebrow != "" {
			b.WriteString("eyebrow: " + p.Eyebrow + "\n")
		}
		b.WriteString("action: " + p.Action + "\n\n" + p.Body + "\n")
		for _, c := range p.Cards {
			b.WriteString("\n## " + c.Title + "\n\nlabel: " + c.Label + "\n")
			if c.Subtitle != "" {
				b.WriteString("subtitle: " + c.Subtitle + "\n")
			}
			if c.Body != "" {
				b.WriteString("\n" + c.Body + "\n")
			}
		}
	}
	return b.String()
}

func chaptersText(chapters []ChapterStory) string {
	var b strings.Builder
	for i, c := range chapters {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("# " + c.Chapter + "\n\n## Briefing\n\n" + c.Briefing + "\n\n## Debrief\n\n" + c.Debrief + "\n")
	}
	return b.String()
}

func termsText(terms []Term) string {
	var b strings.Builder
	for i, t := range terms {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("# " + t.Term + "\n\n")
		if !t.Linked {
			b.WriteString("linked: no\n\n")
		}
		b.WriteString(t.Meaning + "\n")
	}
	return b.String()
}

// ---- Rows ----

// LoadStory reads the story of the course with the given uuid from a library: its pages and
// their cards in order, its chapters' stories in the course's order, and its terms in
// alphabetical order, ignoring case. A course without a story has the empty one.
func LoadStory(q store.Querier, courseUUID string) (Story, error) {
	var s Story
	pages := map[int64]int{} // each page's index in s.Pages, by its position
	err := eachRow(q, `SELECT position, kind, eyebrow, title, body, action FROM story_page WHERE course_uuid = ? ORDER BY position`, courseUUID, func(r *sql.Rows) error {
		var p Page
		var position int64
		var eyebrow sql.NullString
		if err := r.Scan(&position, &p.Kind, &eyebrow, &p.Title, &p.Body, &p.Action); err != nil {
			return err
		}
		p.Eyebrow, pages[position] = eyebrow.String, len(s.Pages)
		s.Pages = append(s.Pages, p)
		return nil
	})
	if err != nil {
		return Story{}, fmt.Errorf("reading story_page: %w", err)
	}
	err = eachRow(q, `SELECT page_position, label, title, subtitle, body FROM story_card WHERE course_uuid = ? ORDER BY page_position, position`, courseUUID, func(r *sql.Rows) error {
		var c Card
		var page int64
		var subtitle, body sql.NullString
		if err := r.Scan(&page, &c.Label, &c.Title, &subtitle, &body); err != nil {
			return err
		}
		i, ok := pages[page]
		if !ok {
			return fmt.Errorf("a card of page %d, which the course does not have", page)
		}
		c.Subtitle, c.Body = subtitle.String, body.String
		s.Pages[i].Cards = append(s.Pages[i].Cards, c)
		return nil
	})
	if err != nil {
		return Story{}, fmt.Errorf("reading story_card: %w", err)
	}
	err = eachRow(q, `SELECT ch.slug, cs.briefing, cs.debrief FROM chapter_story cs JOIN chapter ch ON ch.uuid = cs.chapter_uuid
		WHERE ch.course_uuid = ? ORDER BY ch.position`, courseUUID, func(r *sql.Rows) error {
		var c ChapterStory
		if err := r.Scan(&c.Chapter, &c.Briefing, &c.Debrief); err != nil {
			return err
		}
		s.Chapters = append(s.Chapters, c)
		return nil
	})
	if err != nil {
		return Story{}, fmt.Errorf("reading chapter_story: %w", err)
	}
	err = eachRow(q, `SELECT term, meaning, linked FROM story_term WHERE course_uuid = ?`, courseUUID, func(r *sql.Rows) error {
		var t Term
		var linked int64
		if err := r.Scan(&t.Term, &t.Meaning, &linked); err != nil {
			return err
		}
		t.Linked = linked == 1
		s.Terms = append(s.Terms, t)
		return nil
	})
	if err != nil {
		return Story{}, fmt.Errorf("reading story_term: %w", err)
	}
	slices.SortFunc(s.Terms, func(a, b Term) int {
		if c := compareTerms(a.Term, b.Term); c != 0 {
			return c
		}
		return strings.Compare(a.Term, b.Term)
	})
	return s, nil
}

// eachRow runs a query of one course's rows and calls fn with each row.
func eachRow(q store.Querier, query, courseUUID string, fn func(*sql.Rows) error) error {
	rows, err := q.Query(query, courseUUID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := fn(rows); err != nil {
			return err
		}
	}
	return errors.Join(rows.Err(), rows.Close())
}

// Put makes the story the whole story of the course with the given uuid, in the caller's
// transaction: the course's pages and their cards, its chapters' stories and its terms are
// deleted, and the story's rows inserted, the pages and the cards of each numbered from 0.
// Every chapter the story tells of must be one of the course's.
func (s Story) Put(tx *sql.Tx, courseUUID string) error {
	for _, q := range []string{
		`DELETE FROM story_card WHERE course_uuid = ?`,
		`DELETE FROM story_page WHERE course_uuid = ?`,
		`DELETE FROM chapter_story WHERE chapter_uuid IN (SELECT uuid FROM chapter WHERE course_uuid = ?)`,
		`DELETE FROM story_term WHERE course_uuid = ?`,
	} {
		if _, err := tx.Exec(q, courseUUID); err != nil {
			return fmt.Errorf("clearing the course's story: %w", err)
		}
	}
	for i, p := range s.Pages {
		if _, err := tx.Exec(`INSERT INTO story_page (course_uuid, position, kind, eyebrow, title, body, action) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			courseUUID, i, p.Kind, nullIfEmpty(p.Eyebrow), p.Title, p.Body, p.Action); err != nil {
			return fmt.Errorf("story_page %d: %w", i, err)
		}
		for j, c := range p.Cards {
			if _, err := tx.Exec(`INSERT INTO story_card (course_uuid, page_position, position, label, title, subtitle, body) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				courseUUID, i, j, c.Label, c.Title, nullIfEmpty(c.Subtitle), nullIfEmpty(c.Body)); err != nil {
				return fmt.Errorf("story_card %d of page %d: %w", j, i, err)
			}
		}
	}
	for _, c := range s.Chapters {
		res, err := tx.Exec(`INSERT INTO chapter_story (chapter_uuid, briefing, debrief) SELECT uuid, ?, ? FROM chapter WHERE course_uuid = ? AND slug = ?`,
			c.Briefing, c.Debrief, courseUUID, c.Chapter)
		if err != nil {
			return fmt.Errorf("chapter_story of %s: %w", c.Chapter, err)
		}
		if n, err := res.RowsAffected(); err != nil || n != 1 {
			return errors.Join(err, fmt.Errorf("the course has no chapter %s", c.Chapter))
		}
	}
	for _, t := range s.Terms {
		if _, err := tx.Exec(`INSERT INTO story_term (course_uuid, term, meaning, linked) VALUES (?, ?, ?, ?)`,
			courseUUID, t.Term, t.Meaning, boolInt(t.Linked)); err != nil {
			return fmt.Errorf("story_term %q: %w", t.Term, err)
		}
	}
	return nil
}
