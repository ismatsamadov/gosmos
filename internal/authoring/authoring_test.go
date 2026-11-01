package authoring

// The packages in testdata/packages are invented lessons for a fictional space startup,
// Apogee: one of each type.

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jobloom/gosmos/db"
	"github.com/jobloom/gosmos/internal/store"
)

var examples = []string{"launch-window", "fuel-budget", "mission-quiz"}

func example(name string) string { return filepath.Join("testdata", "packages", name) }

func TestUUIDv5(t *testing.T) {
	dns := mustUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if got := uuidV5(dns, "www.example.com"); got != "2ed6657d-e927-568b-95e1-2665a8aea6a2" {
		t.Errorf("uuidV5(DNS, www.example.com) = %s, want RFC 4122's 2ed6657d-e927-568b-95e1-2665a8aea6a2", got)
	}
	a := LessonUUID("learn-launches", "1-countdown", "launch-window")
	if a != LessonUUID("learn-launches", "1-countdown", "launch-window") || a == LessonUUID("learn-launches", "2-orbit", "launch-window") {
		t.Error("a lesson's uuid must depend on its slugs, and on nothing else")
	}
	if a[14] != '5' {
		t.Errorf("%s is not a version 5 uuid", a)
	}
}

// files reads every file under dir, by path.
func files(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		rel, _ := filepath.Rel(dir, p)
		out[filepath.ToSlash(rel)] = string(b)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func sameFiles(t *testing.T, got, want string) {
	t.Helper()
	g, w := files(t, got), files(t, want)
	for name, content := range w {
		if g[name] != content {
			t.Errorf("%s differs:\n%q\nwant\n%q", name, g[name], content)
		}
	}
	for name := range g {
		if _, ok := w[name]; !ok {
			t.Errorf("%s was not in the package", name)
		}
	}
}

func TestReadWriteRoundTrip(t *testing.T) {
	for _, name := range examples {
		t.Run(name, func(t *testing.T) {
			p, err := Read(example(name))
			if err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(t.TempDir(), name)
			if err := p.Write(out); err != nil {
				t.Fatal(err)
			}
			sameFiles(t, out, example(name))
		})
	}
	p, _ := Read(example("fuel-budget"))
	var flags []string
	for _, f := range p.Starter {
		flags = append(flags, f.Name+map[bool]string{true: " hidden", false: ""}[f.Hidden]+map[bool]string{true: " read-only", false: ""}[f.ReadOnly])
	}
	if want := []string{"fuel.go", "fuel_test.go read-only", "limits.go read-only", "vehicle.go hidden"}; !reflect.DeepEqual(flags, want) {
		t.Errorf("fuel-budget's starter files = %q, want %q", flags, want)
	}
	if len(p.Solution) != 1 || p.Solution[0].Name != "fuel.go" {
		t.Errorf("fuel-budget's solution = %+v", p.Solution)
	}
}

// Into a library and back: Record, then Load, gives the same package and the same files.
func TestRecordAndLoad(t *testing.T) {
	d, err := db.Create(filepath.Join(t.TempDir(), "apogee.sqlite"), db.Info{Title: "Apogee", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, s := range []string{
		`INSERT INTO course VALUES ('c', 'learn-launches', 'Learn Launches', 'Course', 'Launch code.', 'An invented course.', 0)`,
		`INSERT INTO chapter VALUES ('ch', 'c', 0, '1-countdown', 'Countdown', '')`,
	} {
		if _, err := d.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	times := Times{Created: "2026-11-01T10:00:00Z", Modified: "2026-11-01T10:00:00Z"}
	for _, name := range examples {
		p, err := Read(example(name))
		if err != nil {
			t.Fatal(err)
		}
		uuid := LessonUUID("learn-launches", "1-countdown", p.Slug)
		tx, _ := d.Begin()
		if err := store.Put(tx, p.Record(uuid, "ch", times)); err != nil {
			t.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		back, err := Load(d, uuid)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(back, p) {
			t.Errorf("%s came back as\n%+v\nwant\n%+v", name, back, p)
		}
		out := filepath.Join(t.TempDir(), name)
		if err := back.Write(out); err != nil {
			t.Fatal(err)
		}
		sameFiles(t, out, example(name))
	}
	if problems, err := db.Validate(context.Background(), d); err != nil || len(problems) > 0 {
		t.Error(problems, err)
	}
}

// edit copies an example package and changes it: files maps a path to its new content, or ""
// to remove it.
func edit(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	src := example(name)
	err := filepath.WalkDir(src, func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		target := filepath.Join(dir, rel)
		os.MkdirAll(filepath.Dir(target), 0o755)
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		target := filepath.Join(dir, filepath.FromSlash(name))
		if content == "" {
			os.Remove(target)
			continue
		}
		os.MkdirAll(filepath.Dir(target), 0o755)
		if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func front(lines ...string) string {
	return "---\n" + strings.Join(lines, "\n") + "\n---\n# A lesson\n"
}

var goodFront = []string{"slug: launch-window", "title: The Launch Window", "type: code", "position: 0", "difficulty: 10", "topics: constants, printing"}

func with(i int, line string) []string {
	f := append([]string(nil), goodFront...)
	f[i] = line
	return f
}

func TestReadRefuses(t *testing.T) {
	for _, tc := range []struct {
		name    string
		example string
		files   map[string]string
		want    string
	}{
		{"no front matter", "launch-window", map[string]string{"README.md": "# A lesson\n"}, "must start with a front matter block"},
		{"front matter not closed", "launch-window", map[string]string{"README.md": "---\nslug: x\n"}, `no closing "---" line`},
		{"a missing key", "launch-window", map[string]string{"README.md": front(goodFront[:5]...)}, "front matter has no topics"},
		{"an unknown key", "launch-window", map[string]string{"README.md": front(append(goodFront, "author: me")...)}, `front matter line "author: me" is not one of`},
		{"keys out of order", "launch-window", map[string]string{"README.md": front(goodFront[1], goodFront[0], goodFront[2], goodFront[3], goodFront[4], goodFront[5])}, "front matter must read exactly"},
		{"a bad slug", "launch-window", map[string]string{"README.md": front(with(0, "slug: Launch_Window")...)}, `slug "Launch_Window" is not lowercase`},
		{"an unknown type", "launch-window", map[string]string{"README.md": front(with(2, "type: video")...)}, `type "video" is not code, code_tests or quiz`},
		{"a negative position", "launch-window", map[string]string{"README.md": front(with(3, "position: -1")...)}, `position "-1" is not a whole number from 0 up`},
		{"no difficulty", "launch-window", map[string]string{"README.md": front(with(4, "difficulty: 0")...)}, `difficulty "0" is not a whole number from 1 (easiest) to 10`},
		{"a difficulty over 10", "launch-window", map[string]string{"README.md": front(with(4, "difficulty: 11")...)}, `difficulty "11" is not a whole number from 1 (easiest) to 10`},
		{"an old difficulty score", "launch-window", map[string]string{"README.md": front(with(4, "difficulty: 40")...)}, `difficulty "40" is not a whole number from 1 (easiest) to 10`},
		{"topics out of order", "launch-window", map[string]string{"README.md": front(with(5, "topics: printing, constants")...)}, "topics: constants, printing"},
		{"a topic twice", "launch-window", map[string]string{"README.md": front(with(5, "topics: printing, printing")...)}, "topics lists a topic twice"},
		{"blanks after a value", "launch-window", map[string]string{"README.md": front(with(1, "title: The Launch Window ")...)}, "front matter must read exactly"},
		{"a picture in the readme", "launch-window", map[string]string{"README.md": front(goodFront...) + "\n![rocket](https://example.com/rocket.png)\n"}, "README.md holds 1 pictures"},
		{"no expected output", "launch-window", map[string]string{"expected_output.txt": ""}, "a code lesson needs expected_output.txt"},
		{"a solution missing", "launch-window", map[string]string{"solution/main.go": ""}, "solution/ has no main.go, a file the learner edits"},
		{"a solution for a test", "fuel-budget", map[string]string{"solution/fuel_test.go": "package main\n"}, "solution/fuel_test.go solves no file the learner edits"},
		{"a solution for a read-only file", "fuel-budget", map[string]string{"solution/limits.go": "package main\n"}, "solution/limits.go solves no file the learner edits"},
		{"a list naming a missing file", "fuel-budget", map[string]string{"readonly.txt": "limits.go\nnothing.go\n"}, "readonly.txt lists nothing.go, which is not in starter/"},
		{"a list out of order", "fuel-budget", map[string]string{"starter/aaa.go": "package main\n", "readonly.txt": "limits.go\naaa.go\n"}, "readonly.txt must list each file once, in alphabetical order"},
		{"a test in readonly.txt", "fuel-budget", map[string]string{"readonly.txt": "fuel_test.go\nlimits.go\n"}, "readonly.txt lists fuel_test.go, a _test.go file"},
		{"an empty list", "fuel-budget", map[string]string{"hidden.txt": "\n"}, "hidden.txt: every line is one file name"},
		{"tests without a test file", "fuel-budget", map[string]string{"starter/fuel_test.go": "", "solution/fuel_test.go": ""}, "a code_tests lesson needs at least one _test.go file"},
		{"nothing to edit", "launch-window", map[string]string{"readonly.txt": "main.go\n", "solution/main.go": ""}, "starter/ has no file the learner edits"},
		{"an expected output in a code_tests lesson", "fuel-budget", map[string]string{"expected_output.txt": "x\n"}, "expected_output.txt does not belong in a code_tests lesson's package"},
		{"a stray file", "launch-window", map[string]string{"notes.md": "to do\n"}, "notes.md does not belong in a code lesson's package"},
		{"text that is not UTF-8", "launch-window", map[string]string{"starter/logo.bin": "\xff\xfe\x00"}, "logo.bin is not UTF-8 text"},
		{"two right answers", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  > Why A.\n- [x] B\n  > Why B.\n"}, `question "Pick one" has 2 right answers, want exactly one`},
		{"no right answer", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [ ] A\n  > Why not A.\n- [ ] B\n  > Why not B.\n"}, `question "Pick one" has 0 right answers`},
		{"an answer twice", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  > Why A.\n- [ ] A\n  > Why not A.\n"}, `question "Pick one" has the answer "A" twice`},
		{"one answer only", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  > Why A.\n"}, `question "Pick one" has fewer than two answers`},
		{"a quiz not written exactly", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n- [x] A\n  > Why A.\n- [ ] B\n  > Why not B.\n"}, "quiz.md: it must read exactly"},
		{"a blank line before an explanation", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n\n  > Why A.\n- [ ] B\n  > Why not B.\n"}, "quiz.md: it must read exactly"},
		{"an answer before any question", "mission-quiz", map[string]string{"quiz.md": "- [x] A\n"}, "comes before any question"},
		{"an answer without an explanation", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n- [ ] B\n  > Why not B.\n"}, `answer "A" of question "Pick one" has no explanation`},
		{"an explanation before any answer", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n  > Why.\n- [x] A\n  > Why A.\n"}, `explanation "  > Why." follows no answer`},
		{"two explanations of one answer", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  > Why A.\n  > And again.\n- [ ] B\n  > Why not B.\n"}, `answer "A" has two explanation lines`},
		{"an empty explanation", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  >  \n- [ ] B\n  > Why not B.\n"}, `answer "A" has an empty explanation`},
		{"an explanation indented wrong", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n > Why A.\n- [ ] B\n  > Why not B.\n"}, `line " > Why A." is neither`},
		{"a link to a video in an answer", "mission-quiz", map[string]string{"quiz.md": "## Watch?\n\n- [x] [the launch](https://youtu.be/abc) at noon\n  > It is at noon.\n- [ ] No\n  > It is on.\n"}, "question 1, answer 1 holds"},
		{"a picture in an explanation", "mission-quiz", map[string]string{"quiz.md": "## Pick one\n\n- [x] A\n  > See ![the chart](https://example.com/c.png) here.\n- [ ] B\n  > Why not B.\n"}, "question 1, answer 1's explanation holds 1 pictures"},
		{"code files in a quiz", "mission-quiz", map[string]string{"starter/main.go": "package main\n"}, "starter does not belong in a quiz lesson's package"},
		{"hints in a quiz", "mission-quiz", map[string]string{"hints.md": "One.\n\nTwo.\n"}, "hints.md does not belong in a quiz lesson's package"},
		{"no hints", "launch-window", map[string]string{"hints.md": ""}, "a code lesson's package needs hints.md: 2 or 3 hints"},
		{"one hint", "launch-window", map[string]string{"hints.md": "Look at the constant.\n"}, "hints.md has 1 hint; a code lesson has 2 or 3"},
		{"four hints", "fuel-budget", map[string]string{"hints.md": "One.\n\nTwo.\n\nThree.\n\nFour.\n"}, "hints.md has 4 hints; a code lesson has 2 or 3"},
		{"a hint that is a list", "launch-window", map[string]string{"hints.md": "- Look at the constant.\n- Then at Printf.\n\nSecond hint.\n"}, "hints.md: hint 1 is not one paragraph of text"},
		{"a hint with a code block", "launch-window", map[string]string{"hints.md": "First hint.\n\n```go\nfmt.Printf(format, windowOpens)\n```\n"}, "hints.md: hint 2 is not one paragraph of text"},
		{"two blank lines between hints", "launch-window", map[string]string{"hints.md": "First hint.\n\n\nSecond hint.\n"}, "hints.md: hint 2 starts with a blank line"},
		{"a hint line ending in spaces", "launch-window", map[string]string{"hints.md": "First hint.  \n\nSecond hint.\n"}, `hints.md: hint 1: line "First hint.  " starts or ends with a space`},
		{"hints without their last newline", "launch-window", map[string]string{"hints.md": "First hint.\n\nSecond hint."}, "hints.md: it must end with exactly one newline"},
		{"a picture in a hint", "launch-window", map[string]string{"hints.md": "See ![the pad](https://example.com/pad.png) first.\n\nSecond hint.\n"}, "hint 1 holds 1 pictures"},
		{"no solution note", "fuel-budget", map[string]string{"solution_note.md": ""}, "a code lesson's package needs solution_note.md"},
		{"a solution note of two paragraphs", "fuel-budget", map[string]string{"solution_note.md": "It subtracts.\n\nAnd it clamps.\n"}, "solution_note.md has 2 paragraphs; it is one"},
		{"an empty solution note", "fuel-budget", map[string]string{"solution_note.md": "\n"}, "solution_note.md: it is empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Read(edit(t, tc.example, tc.files))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Read: err = %v, want it to say %q", err, tc.want)
			}
		})
	}
}

func TestIgnoresDotFiles(t *testing.T) {
	dir := edit(t, "launch-window", map[string]string{".DS_Store": "x", "starter/.notes": "x"})
	p, err := Read(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Starter) != 1 {
		t.Errorf("starter files %+v, want main.go only", p.Starter)
	}
	if err := p.Write(dir); err == nil || !bytes.Contains([]byte(err.Error()), []byte("not empty")) {
		t.Errorf("Write into a folder that has files: err = %v", err)
	}
}

// indent indents every line that is not empty by two spaces, as Layout shows its examples.
func indent(text string) string {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = "  " + l
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

// The examples the help shows are what Read takes, so writers can start from them.
func TestLayoutExamplesAreValid(t *testing.T) {
	front := "---\nslug: launch-checklist\ntitle: The Launch Checklist\ntype: code\nposition: 0\ndifficulty: 3\ntopics: printing, variables\n---\n"
	hints := "Look at what clock.go defines: the number you need is already there.\n\n" +
		"fmt.Printf puts each value you give it where a %d stands, in order.\n\n" +
		"Compare your output with the expected one character by character, the newline too.\n"
	note := "The constant holds the number, and Printf puts it where the question mark was, followed\nby the newline the expected output ends with.\n"
	quiz := "## Which verb prints a value in Go's default format?\n\n" +
		"- [ ] %d\n  > %d prints integers only; other values come out as an error marker.\n" +
		"- [x] %v\n  > %v prints any value in its default format.\n" +
		"- [ ] %q\n  > %q prints a quoted string, not the default format.\n\n" +
		"## Does fmt.Println add a newline?\n\n" +
		"- [x] Yes\n  > Println always ends what it prints with a newline.\n" +
		"- [ ] No\n  > That is fmt.Print, which adds nothing.\n"
	for name, example := range map[string]string{"front matter": front, "hints": hints, "solution note": note, "quiz": quiz} {
		if !strings.Contains(Layout, indent(example)) {
			t.Errorf("Layout does not show the %s example", name)
		}
	}
	var p Package
	if err := p.parseReadme(front); err != nil {
		t.Errorf("the front matter example: %v", err)
	}
	if got, err := paragraphs(hints, "hint"); err != nil || len(got) != 3 {
		t.Errorf("the hints example: %d hints, %v", len(got), err)
	}
	if got, err := paragraphs(note, "note"); err != nil || len(got) != 1 {
		t.Errorf("the solution note example: %d paragraphs, %v", len(got), err)
	}
	if qs, err := parseQuiz(quiz); err != nil || len(qs) != 2 || qs[0].Answers[1].Explanation != "%v prints any value in its default format." {
		t.Errorf("the quiz example: %+v, %v", qs, err)
	}
}
