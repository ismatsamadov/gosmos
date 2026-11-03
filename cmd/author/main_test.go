package main

// Every course, chapter and lesson here is invented: lessons for a fictional space startup,
// Apogee, from internal/authoring/testdata/packages.

import (
	"bytes"
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jobloom/gosmos/db"
	"github.com/jobloom/gosmos/internal/authoring"
)

func example(name string) string {
	return filepath.Join("..", "..", "internal", "authoring", "testdata", "packages", name)
}

func author(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(context.Background(), args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func mustAuthor(t *testing.T, args ...string) string {
	t.Helper()
	code, out, errOut := author(t, args...)
	if code != 0 {
		t.Fatalf("author %v: exit %d: %s", args, code, errOut)
	}
	return out
}

// library makes a library with the course learn-launches and its chapter 1-countdown.
func library(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "apogee.sqlite")
	d, err := db.Create(path, db.Info{Title: "Apogee", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	d.Close()
	mustAuthor(t, "course", path, "--slug", "learn-launches", "--title", "Learn Launches", "--kind", "Course",
		"--short", "Write the code that flies Apogee's rockets.", "--description", "A course of invented lessons.")
	mustAuthor(t, "chapter", path, "--course", "learn-launches", "--slug", "1-countdown", "--title", "Countdown",
		"--description", "The minutes before liftoff.")
	return path
}

func importAll(t *testing.T, path string) {
	t.Helper()
	for _, name := range []string{"launch-window", "fuel-budget", "mission-quiz"} {
		out := mustAuthor(t, "lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", example(name))
		if !strings.Contains(out, "added") {
			t.Errorf("import of %s: %q", name, out)
		}
	}
}

func open(t *testing.T, path string) *sql.DB {
	t.Helper()
	d, _, err := db.Inspect(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func count(t *testing.T, path, query string, args ...any) int {
	t.Helper()
	var n int
	if err := open(t, path).QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

func sameTree(t *testing.T, got, want string) {
	t.Helper()
	read := func(dir string) map[string]string {
		out := map[string]string{}
		filepath.WalkDir(dir, func(p string, e fs.DirEntry, err error) error {
			if err == nil && !e.IsDir() {
				b, _ := os.ReadFile(p)
				rel, _ := filepath.Rel(dir, p)
				out[filepath.ToSlash(rel)] = string(b)
			}
			return err
		})
		return out
	}
	g, w := read(got), read(want)
	if len(g) != len(w) {
		t.Errorf("%s has %d files, %s has %d", got, len(g), want, len(w))
	}
	for name, content := range w {
		if g[name] != content {
			t.Errorf("%s differs from %s", filepath.Join(got, name), filepath.Join(want, name))
		}
	}
}

func TestAuthorALibrary(t *testing.T) {
	path := library(t)
	importAll(t, path)

	d := open(t, path)
	if problems, err := db.Validate(context.Background(), d); err != nil || len(problems) > 0 {
		t.Fatal(problems, err)
	}
	for query, want := range map[string]int{
		`SELECT count(*) FROM lesson`:                                        3,
		`SELECT count(*) FROM lesson WHERE uuid IN (?, ?, ?)`:                3,
		`SELECT count(*) FROM starter_file`:                                  5,
		`SELECT count(*) FROM starter_file WHERE is_hidden = 1`:              1,
		`SELECT count(*) FROM starter_file WHERE is_readonly = 1`:            2,
		`SELECT count(*) FROM solution_file`:                                 2,
		`SELECT count(*) FROM code_lesson WHERE expected_output IS NOT NULL`: 1,
		`SELECT count(*) FROM quiz_answer WHERE is_correct = 1`:              2,
		`SELECT count(*) FROM lesson_topic`:                                  4,
		`SELECT count(*) FROM lesson WHERE created_at = last_mod`:            3,
		`SELECT count(*) FROM lesson WHERE difficulty BETWEEN 1 AND 3`:       3,
	} {
		var args []any
		if strings.Contains(query, "IN (?") {
			for _, s := range []string{"launch-window", "fuel-budget", "mission-quiz"} {
				args = append(args, authoring.LessonUUID("learn-launches", "1-countdown", s))
			}
		}
		var n int
		if err := d.QueryRow(query, args...).Scan(&n); err != nil || n != want {
			t.Errorf("%s = %d, want %d (%v)", query, n, want, err)
		}
	}

	// Importing the same packages again writes nothing.
	before, _ := os.ReadFile(path)
	for _, name := range []string{"launch-window", "fuel-budget", "mission-quiz"} {
		if out := mustAuthor(t, "lesson", path, "--chapter", "1-countdown", example(name), "--course", "learn-launches"); !strings.Contains(out, "unchanged") {
			t.Errorf("a second import of %s: %q", name, out)
		}
	}
	if code, out, _ := author(t, "course", path, "--slug", "learn-launches", "--title", "Learn Launches", "--kind", "Course",
		"--short", "Write the code that flies Apogee's rockets.", "--description", "A course of invented lessons."); code != 0 || !strings.Contains(out, "unchanged") {
		t.Errorf("the same course again: %d %q", code, out)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("importing the same lessons again changed the library")
	}

	// Exporting gives back the packages, byte for byte, and importing them changes nothing.
	out := filepath.Join(t.TempDir(), "export")
	if msg := mustAuthor(t, "export", path, "--course", "learn-launches", out); !strings.Contains(msg, "exported 3 lessons") {
		t.Errorf("export: %q", msg)
	}
	for _, name := range []string{"launch-window", "fuel-budget", "mission-quiz"} {
		sameTree(t, filepath.Join(out, "1-countdown", name), example(name))
		if msg := mustAuthor(t, "lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", filepath.Join(out, "1-countdown", name)); !strings.Contains(msg, "unchanged") {
			t.Errorf("re-import of the export of %s: %q", name, msg)
		}
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("re-importing the export changed the library")
	}
	one := filepath.Join(t.TempDir(), "chapter")
	mustAuthor(t, "export", path, "--course", "learn-launches", "--chapter", "1-countdown", one)
	sameTree(t, filepath.Join(one, "fuel-budget"), example("fuel-budget"))
	if code, _, errOut := author(t, "export", path, "--course", "learn-launches", one); code != 1 || !strings.Contains(errOut, "is not empty") {
		t.Errorf("export into a folder with files: %d %q", code, errOut)
	}
}

// copyPackage copies an example package so a test can change it.
func copyPackage(t *testing.T, name string) string {
	t.Helper()
	p, err := authoring.Read(example(name))
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), name)
	if err := p.Write(dir); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestReplaceInPlace(t *testing.T) {
	path := library(t)
	importAll(t, path)
	uuid := authoring.LessonUUID("learn-launches", "1-countdown", "launch-window")
	var created string
	open(t, path).QueryRow(`SELECT created_at FROM lesson WHERE uuid = ?`, uuid).Scan(&created)

	dir := copyPackage(t, "launch-window")
	readme := filepath.Join(dir, "README.md")
	b, _ := os.ReadFile(readme)
	os.WriteFile(readme, append(b, "\nA new closing line.\n"...), 0o644)
	os.WriteFile(filepath.Join(dir, "solution", "main.go"), []byte("package main\n\nimport \"fmt\"\n\nconst windowOpens = 90\n\nfunc main() {\n\tfmt.Println(\"Window opens at T-minus\", windowOpens, \"seconds\")\n}\n"), 0o644)
	out := mustAuthor(t, "lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", dir)
	if !strings.Contains(out, "updated readme, solution_file ("+uuid+")") {
		t.Errorf("import of a changed lesson: %q", out)
	}
	var gotCreated, gotReadme string
	open(t, path).QueryRow(`SELECT created_at, readme FROM lesson WHERE uuid = ?`, uuid).Scan(&gotCreated, &gotReadme)
	if gotCreated != created || !strings.HasSuffix(gotReadme, "A new closing line.\n") {
		t.Errorf("after replacing: created_at %s (was %s), readme ends %q", gotCreated, created, gotReadme[max(0, len(gotReadme)-30):])
	}
	if n := count(t, path, `SELECT count(*) FROM lesson`); n != 3 {
		t.Errorf("%d lessons after replacing one, want 3", n)
	}

	// A lesson made elsewhere, under its own uuid, keeps it when a package replaces it.
	other := library(t)
	exec := func(stmts ...string) {
		d, err := db.Open(other)
		if err != nil {
			t.Fatal(err)
		}
		defer d.Close()
		for _, s := range stmts {
			if _, err := d.Exec(s); err != nil {
				t.Fatalf("%s: %v", s, err)
			}
		}
	}
	chapter := authoring.ChapterUUID("learn-launches", "1-countdown")
	exec(`INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
			VALUES ('made-elsewhere', '`+chapter+`', 0, 'mission-quiz', 'Old quiz', 'type_choice', 1, 5, 'Old.', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`,
		`INSERT INTO quiz_lesson VALUES ('made-elsewhere', 'type_choice', 0)`,
		`INSERT INTO quiz_question VALUES ('made-elsewhere', 0, 'Old?')`,
		`INSERT INTO quiz_answer VALUES ('made-elsewhere', 0, 0, 'Yes', 1, NULL), ('made-elsewhere', 0, 1, 'No', 0, NULL)`)
	quiz := copyPackage(t, "mission-quiz")
	front := filepath.Join(quiz, "README.md")
	b, _ = os.ReadFile(front)
	os.WriteFile(front, bytes.Replace(b, []byte("position: 2"), []byte("position: 0"), 1), 0o644)
	if out := mustAuthor(t, "lesson", other, "--course", "learn-launches", "--chapter", "1-countdown", quiz); !strings.Contains(out, "(made-elsewhere)") {
		t.Errorf("replacing a lesson made elsewhere: %q", out)
	}
	if n := count(t, other, `SELECT count(*) FROM lesson WHERE uuid = 'made-elsewhere' AND title = 'Mission Check'`); n != 1 {
		t.Error("the lesson made elsewhere was not replaced in place, under its own uuid")
	}
	if n := count(t, other, `SELECT count(*) FROM quiz_answer WHERE lesson_uuid = 'made-elsewhere' AND explanation IS NOT NULL`); n != 5 {
		t.Errorf("%d answers of the replaced quiz have explanations, want 5", n)
	}
}

func TestAuthorRefuses(t *testing.T) {
	path := library(t)
	importAll(t, path)
	taken := copyPackage(t, "launch-window")
	front := filepath.Join(taken, "README.md")
	b, _ := os.ReadFile(front)
	os.WriteFile(front, bytes.Replace(bytes.Replace(b, []byte("launch-window"), []byte("second-window"), 1), []byte("position: 0"), []byte("position: 1"), 1), 0o644)
	bad := copyPackage(t, "mission-quiz")
	os.WriteFile(filepath.Join(bad, "quiz.md"), []byte("## Pick one\n\n- [x] A\n  > Why A.\n- [x] B\n  > Why B.\n"), 0o644)
	unexplained := copyPackage(t, "mission-quiz")
	os.WriteFile(filepath.Join(unexplained, "quiz.md"), []byte("## Pick one\n\n- [x] A\n  > Why A.\n- [ ] B\n"), 0o644)
	hintless := copyPackage(t, "launch-window")
	os.Remove(filepath.Join(hintless, "hints.md"))
	oneHint := copyPackage(t, "launch-window")
	os.WriteFile(filepath.Join(oneHint, "hints.md"), []byte("Look at the constant.\n"), 0o644)
	noteless := copyPackage(t, "fuel-budget")
	os.Remove(filepath.Join(noteless, "solution_note.md"))
	oldScore := copyPackage(t, "fuel-budget")
	readme, _ := os.ReadFile(filepath.Join(oldScore, "README.md"))
	os.WriteFile(filepath.Join(oldScore, "README.md"), bytes.Replace(readme, []byte("difficulty: 3"), []byte("difficulty: 40"), 1), 0o644)
	before, _ := os.ReadFile(path)

	for _, tc := range []struct {
		name string
		args []string
		code int
		want string
	}{
		{"a position another lesson has", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", taken}, 1, "position 1 of chapter 1-countdown is taken by lesson fuel-budget"},
		{"a broken package", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", bad}, 1, `has 2 right answers, want exactly one`},
		{"a quiz answer without an explanation", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", unexplained}, 1, `answer "B" of question "Pick one" has no explanation`},
		{"a code lesson without hints", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", hintless}, 1, "a code lesson's package needs hints.md: 2 or 3 hints"},
		{"a code lesson with one hint", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", oneHint}, 1, "hints.md has 1 hint; a code lesson has 2 or 3"},
		{"a code lesson without a solution note", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", noteless}, 1, "a code lesson's package needs solution_note.md"},
		{"a difficulty off the scale", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", oldScore}, 1, `difficulty "40" is not a whole number from 1 (easiest) to 10`},
		{"an unknown course", []string{"lesson", path, "--course", "learn-landings", "--chapter", "1-countdown", example("launch-window")}, 1, "the library has no course learn-landings"},
		{"an unknown chapter", []string{"lesson", path, "--course", "learn-launches", "--chapter", "9-orbit", example("launch-window")}, 1, "course learn-launches has no chapter 9-orbit"},
		{"no chapter flag", []string{"lesson", path, "--course", "learn-launches", example("launch-window")}, 2, "--course and --chapter are required"},
		{"an unknown kind", []string{"course", path, "--slug", "learn-orbits", "--title", "Orbits", "--kind", "Lecture", "--short", "S", "--description", "D"}, 2, `--kind "Lecture" is not Course or Guided Project`},
		{"a bad course slug", []string{"course", path, "--slug", "Learn Orbits", "--title", "Orbits", "--kind", "Course", "--short", "S", "--description", "D"}, 2, "is not lowercase letters"},
		{"a picture in a description", []string{"course", path, "--slug", "learn-orbits", "--title", "Orbits", "--kind", "Course", "--short", "S", "--description", "![orbit](https://example.com/orbit.png)"}, 1, "--description holds pictures"},
		{"a course position taken", []string{"course", path, "--slug", "learn-orbits", "--title", "Orbits", "--kind", "Course", "--short", "S", "--description", "D", "--position", "0"}, 1, "position 0 is taken by course learn-launches"},
		{"a chapter of no course", []string{"chapter", path, "--course", "learn-orbits", "--slug", "1-ellipses", "--title", "Ellipses", "--description", "D"}, 1, "the library has no course learn-orbits"},
		{"an unknown command", []string{"publish", path}, 2, `unknown command "publish"`},
		{"too many arguments", []string{"lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", example("launch-window"), "extra"}, 2, "want 2 arguments"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, _, errOut := author(t, tc.args...)
			if code != tc.code || !strings.Contains(errOut, tc.want) {
				t.Errorf("exit %d, stderr %q; want %d and %q", code, errOut, tc.code, tc.want)
			}
		})
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("a refused command changed the library")
	}
}

func storyExample() string {
	return filepath.Join("..", "..", "internal", "authoring", "testdata", "story")
}

// storyWith copies the example story package and sets a file of it: to content, or, for "",
// leaves it out.
func storyWith(t *testing.T, name, content string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "story")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(storyExample())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(storyExample(), e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, e.Name()), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	target := filepath.Join(dir, name)
	if content == "" {
		os.Remove(target)
	} else if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// A course's story goes into the library whole and comes back out byte for byte. Importing it
// again, or its export, writes nothing; a changed file updates the story, and a folder with none
// of the files clears it.
func TestStory(t *testing.T) {
	path := library(t)
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", storyExample()); out != "story learn-launches: added (2 pages, 4 cards, 1 chapter and 4 terms)\n" {
		t.Errorf("story: %q", out)
	}
	for query, want := range map[string]int{
		`SELECT count(*) FROM story_page`: 2,
		`SELECT count(*) FROM story_page WHERE position = 1 AND kind = 'badge' AND eyebrow = 'APOGEE'`: 1,
		`SELECT count(*) FROM story_card`:                                                                                          4,
		`SELECT count(*) FROM story_card WHERE subtitle IS NULL AND body IS NULL`:                                                  2,
		`SELECT count(*) FROM chapter_story WHERE chapter_uuid = '` + authoring.ChapterUUID("learn-launches", "1-countdown") + `'`: 1,
		`SELECT count(*) FROM story_term`:                                                                                          4,
		`SELECT count(*) FROM story_term WHERE linked = 0 AND term = 'pad'`:                                                        1,
	} {
		if n := count(t, path, query); n != want {
			t.Errorf("%s = %d, want %d", query, n, want)
		}
	}

	before, _ := os.ReadFile(path)
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", storyExample()); out != "story learn-launches: unchanged (2 pages, 4 cards, 1 chapter and 4 terms)\n" {
		t.Errorf("a second import: %q", out)
	}
	out := filepath.Join(t.TempDir(), "story")
	if msg := mustAuthor(t, "export", path, "--course", "learn-launches", "--story", out); msg != "exported the story of learn-launches to "+out+" (2 pages, 4 cards, 1 chapter and 4 terms)\n" {
		t.Errorf("export --story: %q", msg)
	}
	sameTree(t, out, storyExample())
	if msg := mustAuthor(t, "story", path, out, "--course", "learn-launches"); !strings.Contains(msg, "unchanged") {
		t.Errorf("re-import of the export: %q", msg)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("importing the same story again changed the library")
	}
	if code, _, errOut := author(t, "export", path, "--course", "learn-launches", "--story", out); code != 1 || !strings.Contains(errOut, "is not empty") {
		t.Errorf("export --story into a folder with files: %d %q", code, errOut)
	}

	terms := "# go/no-go\n\nThe poll before a launch: each console says go, or calls a hold.\n\n# T-minus\n\nThe time left before liftoff.\n"
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", storyWith(t, "terms.md", terms)); out != "story learn-launches: updated terms.md (2 pages, 4 cards, 1 chapter and 2 terms)\n" {
		t.Errorf("a story with other terms: %q", out)
	}
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", storyWith(t, "opening.md", "")); out != "story learn-launches: updated opening.md, terms.md (1 chapter and 4 terms)\n" {
		t.Errorf("a story without its opening: %q", out)
	}
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", t.TempDir()); out != "story learn-launches: cleared (it had 1 chapter and 4 terms)\n" {
		t.Errorf("an empty story folder: %q", out)
	}
	if n := count(t, path, `SELECT (SELECT count(*) FROM story_page) + (SELECT count(*) FROM story_card) + (SELECT count(*) FROM chapter_story) + (SELECT count(*) FROM story_term)`); n != 0 {
		t.Errorf("%d story rows are left after clearing the story", n)
	}
	none := filepath.Join(t.TempDir(), "none")
	if msg := mustAuthor(t, "export", path, "--course", "learn-launches", "--story", none); msg != "exported the story of learn-launches to "+none+" (no story)\n" {
		t.Errorf("export --story of no story: %q", msg)
	}
	if entries, err := os.ReadDir(none); err != nil || len(entries) != 0 {
		t.Errorf("the export of no story holds %v, %v; want an empty folder", entries, err)
	}
	if out := mustAuthor(t, "story", path, "--course", "learn-launches", none); out != "story learn-launches: unchanged (no story)\n" {
		t.Errorf("re-import of no story: %q", out)
	}
	if problems, err := db.Validate(context.Background(), open(t, path)); err != nil || len(problems) > 0 {
		t.Error(problems, err)
	}
}

func TestStoryRefuses(t *testing.T) {
	path := library(t)
	mustAuthor(t, "chapter", path, "--course", "learn-launches", "--slug", "2-orbit", "--title", "Orbit", "--description", "Around and around.")
	mustAuthor(t, "story", path, "--course", "learn-launches", storyExample())
	story := "\n\n## Briefing\n\nThe countdown starts.\n\n## Debrief\n\nIt ended well.\n"
	unknown := storyWith(t, "chapters.md", "# 9-landing"+story)
	disordered := storyWith(t, "chapters.md", "# 2-orbit"+story+"\n# 1-countdown"+story)
	broken := storyWith(t, "terms.md", "# pad\n\nlinked: yes\n\nThe concrete a rocket stands on.\n")
	stray := storyWith(t, "notes.md", "to do\n")
	before, _ := os.ReadFile(path)

	for _, tc := range []struct {
		name string
		args []string
		code int
		want string
	}{
		{"a chapter the course does not have", []string{"story", path, "--course", "learn-launches", unknown}, 1, "chapters.md line 1: course learn-launches has no chapter 9-landing; add it with author chapter, or leave its story out"},
		{"chapters out of the course's order", []string{"story", path, "--course", "learn-launches", disordered}, 1, "chapters.md line 11: chapter 1-countdown comes after 2-orbit, but the course has it before"},
		{"a broken package", []string{"story", path, "--course", "learn-launches", broken}, 1, `package ` + broken + `: terms.md line 3: "linked: yes" is not "linked: no"`},
		{"a stray file", []string{"story", path, "--course", "learn-launches", stray}, 1, "notes.md does not belong in a story package"},
		{"a missing folder", []string{"story", path, "--course", "learn-launches", filepath.Join(t.TempDir(), "nowhere")}, 1, "no such file or directory"},
		{"an unknown course", []string{"story", path, "--course", "learn-landings", storyExample()}, 1, "the library has no course learn-landings"},
		{"no course flag", []string{"story", path, storyExample()}, 2, "--course is required"},
		{"too many arguments", []string{"story", path, "--course", "learn-launches", storyExample(), "extra"}, 2, "want 2 arguments"},
		{"the story of one chapter", []string{"export", path, "--course", "learn-launches", "--chapter", "1-countdown", "--story", filepath.Join(t.TempDir(), "out")}, 2, "--chapter and --story do not go together"},
		{"the story of an unknown course", []string{"export", path, "--course", "learn-landings", "--story", filepath.Join(t.TempDir(), "out")}, 1, "the library has no course learn-landings"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, _, errOut := author(t, tc.args...)
			if code != tc.code || !strings.Contains(errOut, tc.want) {
				t.Errorf("exit %d, stderr %q; want %d and %q", code, errOut, tc.code, tc.want)
			}
		})
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("a refused command changed the library")
	}
}

// A write that would leave the library failing its checks is not made.
func TestAuthorChecksTheLibrary(t *testing.T) {
	path := library(t)
	d, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec(`INSERT INTO lesson_hint VALUES ('no-such-lesson', 0, 'Lost.')`); err != nil {
		t.Fatal(err)
	}
	d.Close()
	code, _, errOut := author(t, "lesson", path, "--course", "learn-launches", "--chapter", "1-countdown", example("launch-window"))
	if code != 1 || !strings.Contains(errOut, "nothing was written") || !strings.Contains(errOut, "foreign key: lesson_hint") {
		t.Errorf("import into a library that fails its checks: %d %q", code, errOut)
	}
	if n := count(t, path, `SELECT count(*) FROM lesson`); n != 0 {
		t.Error("the refused import left a lesson behind")
	}
	code, _, errOut = author(t, "story", path, "--course", "learn-launches", storyExample())
	if code != 1 || !strings.Contains(errOut, "nothing was written") || !strings.Contains(errOut, "foreign key: lesson_hint") {
		t.Errorf("a story imported into a library that fails its checks: %d %q", code, errOut)
	}
	if n := count(t, path, `SELECT (SELECT count(*) FROM story_page) + (SELECT count(*) FROM story_term)`); n != 0 {
		t.Error("the refused import left a story behind")
	}
}

func TestHelp(t *testing.T) {
	code, out, _ := author(t, "help")
	if code != 0 || !strings.Contains(out, "A lesson package is a folder:") || !strings.Contains(out, "author lesson <file>") ||
		!strings.Contains(out, "hints.md holds 2 or 3 hints") || !strings.Contains(out, "solution_note.md says why") || !strings.Contains(out, `two spaces, "> ", then the`) {
		t.Errorf("help: %d %q", code, out[:min(len(out), 200)])
	}
	if !strings.Contains(out, "author story <file> --course S <story folder>") || !strings.Contains(out, "[--chapter S | --story]") ||
		!strings.Contains(out, "A story package is a folder") || !strings.Contains(out, "opening.md holds one or more pages") ||
		!strings.Contains(out, "chapters.md holds the briefing") || !strings.Contains(out, "terms.md holds the field guide's terms") {
		t.Errorf("help does not describe the story package: %q", out[max(0, len(out)-300):])
	}
}
