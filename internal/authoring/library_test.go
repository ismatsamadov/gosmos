package authoring

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ismatsamadov/gosmos/db"
)

// libraryFolder writes a library folder made of the example packages and story: one course,
// learn-launches, with one chapter, 1-countdown.
func libraryFolder(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "apogee")
	write := func(name, text string) {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("library.txt", "uuid: 0b7f3e52-9a41-4c6d-8f20-6e5d1c9b7a34\ntitle: Apogee\nversion: 1\n")
	write("learn-launches/course.md", "---\nslug: learn-launches\ntitle: Learn Launches\nkind: Course\nposition: 0\nshort: Launch code.\n---\nAn invented course.\n")
	write("learn-launches/1-countdown/chapter.md", "---\nslug: 1-countdown\ntitle: Countdown\nposition: 0\n---\nThe minutes before liftoff.\n")
	for _, name := range []string{"opening.md", "chapters.md", "terms.md"} {
		b, err := os.ReadFile(filepath.Join(storyExample(), name))
		if err != nil {
			t.Fatal(err)
		}
		write("learn-launches/story/"+name, string(b))
	}
	for _, name := range examples {
		p, err := Read(example(name))
		if err != nil {
			t.Fatal(err)
		}
		p.Created, p.Modified = "2026-11-01T10:00:00Z", "2026-11-02T08:30:00Z"
		if err := p.Write(filepath.Join(dir, "learn-launches", "1-countdown", name)); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLibraryRoundTrip(t *testing.T) {
	src := libraryFolder(t)
	l, err := ReadLibrary(src)
	if err != nil {
		t.Fatal(err)
	}
	if got := l.Summary(); got != "1 course, 1 chapter and 3 lessons" {
		t.Errorf("Summary() = %q", got)
	}
	if c := l.Courses[0]; c.Story.Empty() || len(c.Chapters[0].Lessons) != 3 || c.Chapters[0].Lessons[2].Slug != "mission-quiz" {
		t.Errorf("the course read as %+v", c)
	}
	out := filepath.Join(t.TempDir(), "out")
	if err := l.Write(out); err != nil {
		t.Fatal(err)
	}
	sameFiles(t, out, src)
}

// Into a new library and back: the same rows from the same folder, in position order, with
// the uuids the slugs give.
func TestBuildAndLoadLibrary(t *testing.T) {
	l, err := ReadLibrary(libraryFolder(t))
	if err != nil {
		t.Fatal(err)
	}
	d, err := db.Create(filepath.Join(t.TempDir(), "apogee.sqlite"), l.Info)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	tx, err := d.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Put(tx); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if problems, err := db.Validate(context.Background(), d); err != nil || len(problems) > 0 {
		t.Fatal(problems, err)
	}
	back, err := LoadLibrary(d)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, l) {
		t.Errorf("the library came back as\n%+v\nwant\n%+v", back, l)
	}
	rows, err := d.Query(`SELECT uuid FROM lesson ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	var want []string
	for _, name := range examples {
		want = append(want, LessonUUID("learn-launches", "1-countdown", name))
	}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("lessons by id = %q, want the slugs' uuids in position order %q", ids, want)
	}
}

func TestReadLibraryRefuses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(dir string) error
		want   string
	}{
		{"no library.txt", func(dir string) error { return os.Remove(filepath.Join(dir, "library.txt")) }, "library.txt"},
		{"library.txt written loosely", func(dir string) error {
			return os.WriteFile(filepath.Join(dir, "library.txt"), []byte("uuid: 0b7f3e52-9a41-4c6d-8f20-6e5d1c9b7a34\nversion: 1\ntitle: Apogee\n"), 0o644)
		}, "library.txt: it must read exactly"},
		{"a stray file", func(dir string) error {
			return os.WriteFile(filepath.Join(dir, "notes.md"), []byte("to do\n"), 0o644)
		}, "notes.md does not belong in a library folder"},
		{"a course folder named for another slug", func(dir string) error {
			return os.Rename(filepath.Join(dir, "learn-launches"), filepath.Join(dir, "learn-landings"))
		}, "learn-landings/course.md: its slug is learn-launches; the folder is named for it"},
		{"a course of another kind", func(dir string) error {
			return replace(filepath.Join(dir, "learn-launches", "course.md"), "kind: Course", "kind: Lecture")
		}, `kind "Lecture" is not Course or Guided Project`},
		{"a course without a description", func(dir string) error {
			return replace(filepath.Join(dir, "learn-launches", "course.md"), "An invented course.\n", "")
		}, "learn-launches/course.md: its description, after the front matter, is empty"},
		{"a picture in a description", func(dir string) error {
			return replace(filepath.Join(dir, "learn-launches", "1-countdown", "chapter.md"), "The minutes", "![clock](clock.png) The minutes")
		}, "1-countdown/chapter.md's description holds 1 picture"},
		{"a chapter called story", func(dir string) error {
			ch := filepath.Join(dir, "learn-launches", "1-countdown")
			if err := replace(filepath.Join(ch, "chapter.md"), "slug: 1-countdown", "slug: story"); err != nil {
				return err
			}
			if err := os.RemoveAll(filepath.Join(dir, "learn-launches", "story")); err != nil {
				return err
			}
			return os.Rename(ch, filepath.Join(dir, "learn-launches", "story"))
		}, "learn-launches/story/chapter.md: a chapter cannot be called story"},
		{"two chapters at one position", func(dir string) error {
			ch := filepath.Join(dir, "learn-launches", "2-liftoff")
			if err := os.MkdirAll(ch, 0o755); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(ch, "chapter.md"), []byte("---\nslug: 2-liftoff\ntitle: Liftoff\nposition: 0\n---\nThe first seconds.\n"), 0o644)
		}, "chapters 1-countdown and 2-liftoff both have position 0"},
		{"a lesson without its times", func(dir string) error {
			readme := filepath.Join(dir, "learn-launches", "1-countdown", "fuel-budget", "README.md")
			return replace(readme, "created: 2026-11-01T10:00:00Z\nmodified: 2026-11-02T08:30:00Z\n", "")
		}, "1-countdown/fuel-budget/README.md: its front matter has no created: and modified: lines"},
		{"a lesson folder named for another slug", func(dir string) error {
			ch := filepath.Join(dir, "learn-launches", "1-countdown")
			return os.Rename(filepath.Join(ch, "fuel-budget"), filepath.Join(ch, "fuel-plan"))
		}, "1-countdown/fuel-plan/README.md: its slug is fuel-budget"},
		{"two lessons at one position", func(dir string) error {
			return replace(filepath.Join(dir, "learn-launches", "1-countdown", "mission-quiz", "README.md"), "position: 2", "position: 1")
		}, "lessons fuel-budget and mission-quiz both have position 1"},
		{"a briefing for a chapter the course doesn't have", func(dir string) error {
			return replace(filepath.Join(dir, "learn-launches", "story", "chapters.md"), "# 1-countdown", "# 9-orbit")
		}, "learn-launches/story: chapters.md line 1: course learn-launches has no chapter 9-orbit"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := libraryFolder(t)
			if err := tc.change(dir); err != nil {
				t.Fatal(err)
			}
			_, err := ReadLibrary(dir)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("ReadLibrary: err = %v, want it to say %q", err, tc.want)
			}
		})
	}
}

// replace changes the first old in the file to new.
func replace(path, old, new string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(b), old) {
		return os.ErrNotExist
	}
	return os.WriteFile(path, []byte(strings.Replace(string(b), old, new, 1)), 0o644)
}
