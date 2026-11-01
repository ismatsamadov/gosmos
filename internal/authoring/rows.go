package authoring

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/jobloom/gosmos/internal/store"
)

// lessonType maps a package type to the library's lesson type, and back.
var lessonType = map[string]string{"code": "type_code", "code_tests": "type_code_tests", "quiz": "type_choice"}

// Times are the timestamps a lesson's row carries, in RFC 3339 UTC.
type Times struct {
	Created, Modified string
}

// Record returns the rows of the lesson the package makes, with the given uuid, in the given
// chapter.
func (p Package) Record(uuid, chapterUUID string, t Times) store.Record {
	l := store.Lesson{
		UUID: uuid, ChapterUUID: chapterUUID, Position: p.Position, Slug: p.Slug, Title: p.Title,
		Type: lessonType[p.Type], Required: true, Difficulty: p.Difficulty, Readme: p.Readme,
		CreatedAt: t.Created, LastMod: t.Modified,
	}
	var rows []store.Row
	add := func(table string, vals ...any) {
		rows = append(rows, store.Row{Table: table, Vals: append([]any{uuid}, vals...)})
	}
	for _, topic := range p.Topics {
		add("lesson_topic", topic)
	}
	for i, hint := range p.Hints {
		add("lesson_hint", int64(i), hint)
	}
	switch p.Type {
	case "quiz":
		add("quiz_lesson", "type_choice", int64(0))
		for qi, q := range p.Questions {
			add("quiz_question", int64(qi), q.Text)
			for ai, a := range q.Answers {
				add("quiz_answer", int64(qi), int64(ai), a.Text, boolInt(a.Right), nullIfEmpty(a.Explanation))
			}
		}
	default:
		var expected any
		if p.Expected != nil {
			expected = *p.Expected
		}
		add("code_lesson", l.Type, "go", expected, nullIfEmpty(p.SolutionNote))
		for i, f := range p.Starter {
			add("starter_file", f.Name, int64(i), f.Content, boolInt(f.Hidden), boolInt(f.ReadOnly))
		}
		for _, f := range p.Solution {
			add("solution_file", f.Name, f.Content)
		}
	}
	return store.Record{Lesson: l, Rows: rows}
}

// Load reads a lesson of the library as a package.
func Load(q store.Querier, uuid string) (Package, error) {
	l, ok, err := store.Get(q, uuid)
	if err != nil {
		return Package{}, err
	}
	if !ok {
		return Package{}, fmt.Errorf("no lesson %s", uuid)
	}
	p := Package{Slug: l.Slug, Title: l.Title, Position: l.Position, Difficulty: l.Difficulty, Readme: l.Readme}
	for k, v := range lessonType {
		if v == l.Type {
			p.Type = k
		}
	}
	if p.Type == "" {
		return p, fmt.Errorf("lesson %s is a %s lesson; packages hold code, code_tests and quiz lessons only", l.Slug, l.Type)
	}
	rows, err := store.ReadParts(q, uuid)
	if err != nil {
		return p, err
	}
	type answerAt struct {
		at int64
		a  Answer
	}
	starterAt := map[string]int64{}
	answers := map[int64][]answerAt{}
	hints := map[int64]string{}
	for _, r := range rows {
		v := r.Vals
		switch r.Table {
		case "lesson_topic":
			p.Topics = append(p.Topics, v[1].(string))
		case "lesson_hint":
			hints[v[1].(int64)] = v[2].(string)
		case "code_lesson":
			if s, ok := v[3].(string); ok {
				p.Expected = &s
			}
			if s, ok := v[4].(string); ok {
				p.SolutionNote = s
			}
		case "starter_file":
			name := v[1].(string)
			starterAt[name] = v[2].(int64)
			p.Starter = append(p.Starter, File{Name: name, Content: v[3].(string), Hidden: v[4].(int64) == 1, ReadOnly: v[5].(int64) == 1})
		case "solution_file":
			p.Solution = append(p.Solution, File{Name: v[1].(string), Content: v[2].(string)})
		case "quiz_question":
			at := v[1].(int64)
			for int64(len(p.Questions)) <= at {
				p.Questions = append(p.Questions, Question{})
			}
			p.Questions[at].Text = v[2].(string)
		case "quiz_answer":
			qi := v[1].(int64)
			explanation, _ := v[5].(string)
			answers[qi] = append(answers[qi], answerAt{v[2].(int64), Answer{Text: v[3].(string), Right: v[4].(int64) == 1, Explanation: explanation}})
		}
	}
	slices.Sort(p.Topics)
	for _, at := range slices.Sorted(maps.Keys(hints)) {
		p.Hints = append(p.Hints, hints[at])
	}
	slices.SortFunc(p.Starter, func(a, b File) int { return int(starterAt[a.Name] - starterAt[b.Name]) })
	slices.SortFunc(p.Solution, func(a, b File) int { return strings.Compare(a.Name, b.Name) })
	for qi := range p.Questions {
		as := answers[int64(qi)]
		slices.SortFunc(as, func(a, b answerAt) int { return int(a.at - b.at) })
		for _, a := range as {
			p.Questions[qi].Answers = append(p.Questions[qi].Answers, a.a)
		}
	}
	for _, f := range p.Starter {
		if isTest(f.Name) && !f.ReadOnly {
			return p, fmt.Errorf("lesson %s: %s is a _test.go file the learner may edit, which a package cannot say", l.Slug, f.Name)
		}
	}
	for i, f := range p.Starter {
		if i > 0 && p.Starter[i-1].Name > f.Name {
			return p, fmt.Errorf("lesson %s: its starter files are not in name order, which a package cannot say", l.Slug)
		}
	}
	return p, nil
}

// nullIfEmpty is s, or NULL for an empty string.
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func boolInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
