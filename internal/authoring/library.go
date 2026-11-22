package authoring

// A library folder holds a whole library as text: the form this repository keeps its courses
// in (courses/apogee/), from which the author tool builds the library file. LibraryLayout is
// its specification. ReadLibrary reads one and Library.Write writes one; LoadLibrary and
// Library.Put map a library to and from a library file's rows.

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/ismatsamadov/gosmos/db"
	"github.com/ismatsamadov/gosmos/internal/store"
)

// LibraryLayout is the library folder specification, as the author tool's help shows it.
const LibraryLayout = `A library folder holds a whole library as text, one folder a course:

  library.txt                     the library's uuid, title and version
  <course>/course.md              the course
  <course>/story/                 the course's story package, if it has a story
  <course>/<chapter>/chapter.md   a chapter
  <course>/<chapter>/<lesson>/    a lesson package, with its created: and modified: lines

Each folder is named for the slug its file gives. Nothing else belongs in it, and a chapter
cannot be called story. Files and folders whose names start with a dot are ignored, and so is
a folder holding nothing else, such as one git left behind with Finder's .DS_Store in it.

library.txt, every line exactly as shown:

  uuid: 86c2c429-0f1a-4f68-a4d0-9b9cb74aff52
  title: Apogee
  version: 1

course.md, the course's front matter, then its description, one line:

  ---
  slug: go-foundations
  title: Go Foundations
  kind: Course
  position: 0
  short: Learn Go from your first program to goroutines, as a new engineer at Apogee.
  ---
  Your first week at Apogee, ...

kind is Course or Guided Project; short is the short description. chapter.md is the same
for a chapter, without kind and short:

  ---
  slug: hello-ground-station
  title: Hello, Ground Station
  position: 0
  ---
  Write, build and run your first Go programs, ...

Positions count from 0 among a library's courses, a course's chapters and a chapter's
lessons, each once. The uuids of courses, chapters and lessons come from their slugs, so a
library built again keeps every learner's progress.`

// Library is a whole library: its library_info row and its courses.
type Library struct {
	Info    db.Info
	Courses []Course // in position order
}

// Course is a course with its chapters and its story.
type Course struct {
	Slug, Title, Kind, Short, Description string
	Position                              int64
	Chapters                              []Chapter // in position order
	Story                                 Story     // empty for a course without a story
}

// Chapter is a chapter with its lessons.
type Chapter struct {
	Slug, Title, Description string
	Position                 int64
	Lessons                  []Package // in position order
}

// Summary counts the library's courses, chapters and lessons.
func (l Library) Summary() string {
	chapters, lessons := 0, 0
	for _, c := range l.Courses {
		chapters += len(c.Chapters)
		for _, ch := range c.Chapters {
			lessons += len(ch.Lessons)
		}
	}
	return fmt.Sprintf("%s, %s and %s", count(len(l.Courses), "course"), count(chapters, "chapter"), count(lessons, "lesson"))
}

// ---- Reading ----

var reUUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ReadLibrary reads and checks the library folder dir. Every refusal names the file.
func ReadLibrary(dir string) (Library, error) {
	var l Library
	text, err := readText(filepath.Join(dir, "library.txt"))
	if err != nil {
		return l, err
	}
	if l.Info, err = parseInfo(text); err != nil {
		return l, fmt.Errorf("library.txt: %w", err)
	}
	entries, err := visible(dir)
	if err != nil {
		return l, err
	}
	for _, e := range entries {
		if e.Name() == "library.txt" {
			continue
		}
		if !e.IsDir() {
			return l, fmt.Errorf("%s does not belong in a library folder: it holds library.txt and a folder for each course", e.Name())
		}
		c, err := readCourse(filepath.Join(dir, e.Name()))
		if err != nil {
			return l, err
		}
		l.Courses = append(l.Courses, c)
	}
	slices.SortFunc(l.Courses, func(a, b Course) int { return int(a.Position - b.Position) })
	for i := 1; i < len(l.Courses); i++ {
		if l.Courses[i].Position == l.Courses[i-1].Position {
			return l, fmt.Errorf("courses %s and %s both have position %d", l.Courses[i-1].Slug, l.Courses[i].Slug, l.Courses[i].Position)
		}
	}
	return l, nil
}

// visible lists dir's entries, leaving out names that start with a dot and folders holding
// nothing else: git leaves a removed lesson's folder behind when Finder put a .DS_Store in it.
func visible(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var kept []os.DirEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			if dots, err := onlyDotFiles(filepath.Join(dir, e.Name())); err != nil {
				return nil, err
			} else if dots {
				continue
			}
		}
		kept = append(kept, e)
	}
	return kept, nil
}

// onlyDotFiles reports whether every file in dir, at any depth, has a name that starts with a
// dot.
func onlyDotFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		switch {
		case strings.HasPrefix(e.Name(), "."):
		case !e.IsDir():
			return false, nil
		default:
			if dots, err := onlyDotFiles(filepath.Join(dir, e.Name())); err != nil || !dots {
				return false, err
			}
		}
	}
	return true, nil
}

func readCourse(dir string) (Course, error) {
	var c Course
	name := filepath.Base(dir)
	text, err := readText(filepath.Join(dir, "course.md"))
	if err != nil {
		return c, fmt.Errorf("%s: %w", name, err)
	}
	f, err := parseRow(text, []string{"slug", "title", "kind", "position", "short"})
	if err != nil {
		return c, fmt.Errorf("%s/course.md: %w", name, err)
	}
	c.Slug, c.Title, c.Kind, c.Short, c.Description, c.Position = f.values["slug"], f.values["title"], f.values["kind"], f.values["short"], f.body, f.position
	if c.Kind != "Course" && c.Kind != "Guided Project" {
		return c, fmt.Errorf("%s/course.md: kind %q is not Course or Guided Project", name, c.Kind)
	}
	if c.Slug != name {
		return c, fmt.Errorf("%s/course.md: its slug is %s; the folder is named for it", name, c.Slug)
	}
	for what, text := range map[string]string{"title": c.Title, "short": c.Short, "description": c.Description} {
		if err := mediaError(name+"/course.md's "+what, text); err != nil {
			return c, err
		}
	}
	entries, err := visible(dir)
	if err != nil {
		return c, err
	}
	for _, e := range entries {
		switch {
		case e.Name() == "course.md":
		case e.Name() == "story" && e.IsDir():
			if _, err := os.Stat(filepath.Join(dir, "story", "chapter.md")); err == nil {
				return c, fmt.Errorf("%s/story/chapter.md: a chapter cannot be called story, the name of a course's story folder", name)
			}
			if c.Story, err = ReadStory(filepath.Join(dir, "story")); err != nil {
				return c, fmt.Errorf("%s/story: %w", name, err)
			}
		case e.IsDir():
			ch, err := readChapter(filepath.Join(dir, e.Name()))
			if err != nil {
				return c, fmt.Errorf("%s/%w", name, err)
			}
			c.Chapters = append(c.Chapters, ch)
		default:
			return c, fmt.Errorf("%s/%s does not belong in a course's folder: it holds course.md, story/ and a folder for each chapter", name, e.Name())
		}
	}
	slices.SortFunc(c.Chapters, func(a, b Chapter) int { return int(a.Position - b.Position) })
	slugs := make([]string, len(c.Chapters))
	for i, ch := range c.Chapters {
		slugs[i] = ch.Slug
		if i > 0 && ch.Position == c.Chapters[i-1].Position {
			return c, fmt.Errorf("%s: chapters %s and %s both have position %d", name, c.Chapters[i-1].Slug, ch.Slug, ch.Position)
		}
	}
	if err := c.Story.CheckChapters(c.Slug, slugs); err != nil {
		return c, fmt.Errorf("%s/story: %w", name, err)
	}
	return c, nil
}

func readChapter(dir string) (Chapter, error) {
	var ch Chapter
	name := filepath.Base(dir)
	text, err := readText(filepath.Join(dir, "chapter.md"))
	if err != nil {
		return ch, fmt.Errorf("%s: %w", name, err)
	}
	f, err := parseRow(text, []string{"slug", "title", "position"})
	if err != nil {
		return ch, fmt.Errorf("%s/chapter.md: %w", name, err)
	}
	ch.Slug, ch.Title, ch.Description, ch.Position = f.values["slug"], f.values["title"], f.body, f.position
	if ch.Slug != name {
		return ch, fmt.Errorf("%s/chapter.md: its slug is %s; the folder is named for it", name, ch.Slug)
	}
	for what, text := range map[string]string{"title": ch.Title, "description": ch.Description} {
		if err := mediaError(name+"/chapter.md's "+what, text); err != nil {
			return ch, err
		}
	}
	entries, err := visible(dir)
	if err != nil {
		return ch, err
	}
	for _, e := range entries {
		if e.Name() == "chapter.md" {
			continue
		}
		if !e.IsDir() {
			return ch, fmt.Errorf("%s/%s does not belong in a chapter's folder: it holds chapter.md and a folder for each lesson", name, e.Name())
		}
		p, err := Read(filepath.Join(dir, e.Name()))
		if err != nil {
			return ch, fmt.Errorf("%s/%s: %w", name, e.Name(), err)
		}
		switch {
		case p.Slug != e.Name():
			return ch, fmt.Errorf("%s/%s/README.md: its slug is %s; the folder is named for it", name, e.Name(), p.Slug)
		case p.Created == "":
			return ch, fmt.Errorf("%s/%s/README.md: its front matter has no created: and modified: lines, which a lesson in a library folder has", name, e.Name())
		}
		ch.Lessons = append(ch.Lessons, p)
	}
	slices.SortFunc(ch.Lessons, func(a, b Package) int { return int(a.Position - b.Position) })
	for i := 1; i < len(ch.Lessons); i++ {
		if ch.Lessons[i].Position == ch.Lessons[i-1].Position {
			return ch, fmt.Errorf("%s: lessons %s and %s both have position %d", name, ch.Lessons[i-1].Slug, ch.Lessons[i].Slug, ch.Lessons[i].Position)
		}
	}
	return ch, nil
}

func parseInfo(text string) (db.Info, error) {
	var i db.Info
	values := map[string]string{}
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		key, value, _ := strings.Cut(line, ": ")
		values[key] = value
	}
	i.UUID, i.Title, i.Version = values["uuid"], values["title"], values["version"]
	switch {
	case !reUUID.MatchString(i.UUID):
		return i, fmt.Errorf("uuid %q is not a uuid in lowercase hex", i.UUID)
	case i.Title == "" || strings.TrimSpace(i.Title) != i.Title:
		return i, errors.New("title is empty, or starts or ends with a space")
	case i.Version == "" || strings.TrimSpace(i.Version) != i.Version:
		return i, errors.New("version is empty, or starts or ends with a space")
	}
	if want := infoText(i); text != want {
		return i, fmt.Errorf("it must read exactly:\n%s", want)
	}
	return i, nil
}

func infoText(i db.Info) string {
	return "uuid: " + i.UUID + "\ntitle: " + i.Title + "\nversion: " + i.Version + "\n"
}

// row is a course.md or chapter.md as read: its front matter's values, its position, and its
// body, the description.
type row struct {
	values   map[string]string
	position int64
	body     string
}

// parseRow reads a course.md or chapter.md whose front matter holds keys, in that order.
func parseRow(text string, keys []string) (row, error) {
	r := row{values: map[string]string{}}
	if !strings.HasPrefix(text, "---\n") {
		return r, errors.New(`it must start with a front matter block: a "---" line, the keys, and another "---" line`)
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return r, errors.New(`its front matter block has no closing "---" line`)
	}
	for _, line := range strings.Split(text[4:4+end], "\n") {
		key, value, _ := strings.Cut(line, ":")
		if !slices.Contains(keys, key) {
			return r, fmt.Errorf("front matter line %q is not one of %s followed by a colon", line, strings.Join(keys, ", "))
		}
		r.values[key] = strings.TrimSpace(value)
	}
	body := text[4+end+5:]
	r.body = strings.TrimSuffix(body, "\n")
	var err error
	if r.position, err = strconv.ParseInt(r.values["position"], 10, 64); err != nil || r.position < 0 {
		return r, fmt.Errorf("position %q is not a whole number from 0 up", r.values["position"])
	}
	switch {
	case !ValidSlug(r.values["slug"]):
		return r, fmt.Errorf("slug %q is not lowercase letters, digits and single hyphens", r.values["slug"])
	case r.body == "" || strings.HasSuffix(r.body, "\n") || strings.TrimSpace(r.body) != r.body:
		return r, errors.New("its description, after the front matter, is empty, or starts or ends with a space or a blank line")
	}
	for _, k := range keys {
		if r.values[k] == "" {
			return r, fmt.Errorf("front matter has no %s", k)
		}
	}
	if want := rowText(keys, r.values, r.body); text != want {
		return r, fmt.Errorf("it must read exactly:\n%s", want)
	}
	return r, nil
}

func rowText(keys []string, values map[string]string, body string) string {
	var b strings.Builder
	b.WriteString("---\n")
	for _, k := range keys {
		b.WriteString(k + ": " + values[k] + "\n")
	}
	b.WriteString("---\n" + body + "\n")
	return b.String()
}

// ---- Writing ----

// Write writes the library into dir as a library folder; dir must not exist yet or be empty.
func (l Library) Write(dir string) error {
	if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
		return fmt.Errorf("%s is not empty", dir)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "library.txt"), []byte(infoText(l.Info)), 0o644); err != nil {
		return err
	}
	for _, c := range l.Courses {
		cdir := filepath.Join(dir, c.Slug)
		if err := os.MkdirAll(cdir, 0o755); err != nil {
			return err
		}
		text := rowText([]string{"slug", "title", "kind", "position", "short"}, map[string]string{
			"slug": c.Slug, "title": c.Title, "kind": c.Kind, "position": strconv.FormatInt(c.Position, 10), "short": c.Short,
		}, c.Description)
		if err := os.WriteFile(filepath.Join(cdir, "course.md"), []byte(text), 0o644); err != nil {
			return err
		}
		if !c.Story.Empty() {
			if err := c.Story.Write(filepath.Join(cdir, "story")); err != nil {
				return fmt.Errorf("%s's story: %w", c.Slug, err)
			}
		}
		for _, ch := range c.Chapters {
			chdir := filepath.Join(cdir, ch.Slug)
			if err := os.MkdirAll(chdir, 0o755); err != nil {
				return err
			}
			text := rowText([]string{"slug", "title", "position"}, map[string]string{
				"slug": ch.Slug, "title": ch.Title, "position": strconv.FormatInt(ch.Position, 10),
			}, ch.Description)
			if err := os.WriteFile(filepath.Join(chdir, "chapter.md"), []byte(text), 0o644); err != nil {
				return err
			}
			for _, p := range ch.Lessons {
				if err := p.Write(filepath.Join(chdir, p.Slug)); err != nil {
					return fmt.Errorf("%s/%s/%s: %w", c.Slug, ch.Slug, p.Slug, err)
				}
			}
		}
	}
	return nil
}

// ---- Rows ----

// LoadLibrary reads a whole library file as a library folder would hold it.
func LoadLibrary(q store.Querier) (Library, error) {
	var l Library
	if err := q.QueryRow(`SELECT uuid, title, version FROM library_info`).Scan(&l.Info.UUID, &l.Info.Title, &l.Info.Version); err != nil {
		return l, fmt.Errorf("library_info: %w", err)
	}
	var uuids []string
	err := eachRow(q, `SELECT uuid, slug, title, kind, short_description, description, position FROM course ORDER BY position`, func(r *sql.Rows) error {
		var c Course
		var uuid string
		if err := r.Scan(&uuid, &c.Slug, &c.Title, &c.Kind, &c.Short, &c.Description, &c.Position); err != nil {
			return err
		}
		uuids = append(uuids, uuid)
		l.Courses = append(l.Courses, c)
		return nil
	})
	if err != nil {
		return l, err
	}
	for i := range l.Courses {
		c := &l.Courses[i]
		if c.Story, err = LoadStory(q, uuids[i]); err != nil {
			return l, fmt.Errorf("%s's story: %w", c.Slug, err)
		}
		if c.Chapters, err = loadChapters(q, uuids[i]); err != nil {
			return l, fmt.Errorf("%s: %w", c.Slug, err)
		}
	}
	return l, nil
}

func loadChapters(q store.Querier, courseUUID string) ([]Chapter, error) {
	var chapters []Chapter
	var uuids []string
	err := eachRow(q, `SELECT uuid, slug, title, description, position FROM chapter WHERE course_uuid = ? ORDER BY position`, func(r *sql.Rows) error {
		var ch Chapter
		var uuid string
		if err := r.Scan(&uuid, &ch.Slug, &ch.Title, &ch.Description, &ch.Position); err != nil {
			return err
		}
		chapters = append(chapters, ch)
		uuids = append(uuids, uuid)
		return nil
	}, courseUUID)
	if err != nil {
		return nil, err
	}
	for i := range chapters {
		var ids []string
		err := eachRow(q, `SELECT uuid FROM lesson WHERE chapter_uuid = ? ORDER BY position`, func(r *sql.Rows) error {
			var id string
			if err := r.Scan(&id); err != nil {
				return err
			}
			ids = append(ids, id)
			return nil
		}, uuids[i])
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			p, err := Load(q, id)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", chapters[i].Slug, err)
			}
			chapters[i].Lessons = append(chapters[i].Lessons, p)
		}
	}
	return chapters, nil
}

// Put writes the library's courses, chapters, lessons and stories into a new, empty library,
// in that order and each in position order, so the same folder always makes the same rows:
// lesson ids, the search index's keys, follow the course, chapter and lesson positions. The
// uuids come from the slugs.
func (l Library) Put(tx *sql.Tx) error {
	for _, c := range l.Courses {
		courseUUID := CourseUUID(c.Slug)
		if _, err := tx.Exec(`INSERT INTO course (uuid, slug, title, kind, short_description, description, position) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			courseUUID, c.Slug, c.Title, c.Kind, c.Short, c.Description, c.Position); err != nil {
			return fmt.Errorf("course %s: %w", c.Slug, err)
		}
		for _, ch := range c.Chapters {
			chapterUUID := ChapterUUID(c.Slug, ch.Slug)
			if _, err := tx.Exec(`INSERT INTO chapter (uuid, course_uuid, position, slug, title, description) VALUES (?, ?, ?, ?, ?, ?)`,
				chapterUUID, courseUUID, ch.Position, ch.Slug, ch.Title, ch.Description); err != nil {
				return fmt.Errorf("chapter %s/%s: %w", c.Slug, ch.Slug, err)
			}
			for _, p := range ch.Lessons {
				rec := p.Record(LessonUUID(c.Slug, ch.Slug, p.Slug), chapterUUID, Times{Created: p.Created, Modified: p.Modified})
				if err := store.Put(tx, rec); err != nil {
					return fmt.Errorf("lesson %s/%s/%s: %w", c.Slug, ch.Slug, p.Slug, err)
				}
			}
		}
	}
	for _, c := range l.Courses {
		if !c.Story.Empty() {
			if err := c.Story.Put(tx, CourseUUID(c.Slug)); err != nil {
				return fmt.Errorf("%s's story: %w", c.Slug, err)
			}
		}
	}
	return nil
}
