package main

// Every lesson here is invented: small programs about a fictional space startup.

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"
)

func q(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// codeLessonSQL returns the statements that add a code lesson in chapter ch-verify: its starter
// files (name -> content; names ending in _test.go are read-only) and its solutions, with two
// hints and a solution note.
func codeLessonSQL(slug, typ string, position int, expected *string, starter, solution map[string]string) []string {
	uuid := "l-" + slug
	exp := "NULL"
	if expected != nil {
		exp = q(*expected)
	}
	stmts := []string{
		`INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
			VALUES ('` + uuid + `', 'ch-verify', ` + fmt.Sprint(position) + `, '` + slug + `', 'Lesson ` + slug + `', '` + typ + `', 1, 10, 'Read me.',
			'2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', '2026-01-02T00:00:00Z', '2026-01-03T00:00:00Z')`,
		`INSERT INTO code_lesson VALUES ('` + uuid + `', '` + typ + `', 'go', ` + exp + `, 'It does what the lesson asks.')`,
		`INSERT INTO lesson_hint VALUES ('` + uuid + `', 0, 'Read the assignment again.'), ('` + uuid + `', 1, 'Look at what the starter does now.')`,
	}
	i := 0
	for _, name := range sortedNames(starter) {
		readonly := 0
		if strings.HasSuffix(name, "_test.go") {
			readonly = 1
		}
		stmts = append(stmts, fmt.Sprintf(`INSERT INTO starter_file VALUES ('%s', %s, %d, %s, 0, %d)`, uuid, q(name), i, q(starter[name]), readonly))
		i++
	}
	for _, name := range sortedNames(solution) {
		stmts = append(stmts, fmt.Sprintf(`INSERT INTO solution_file VALUES ('%s', %s, %s)`, uuid, q(name), q(solution[name])))
	}
	return stmts
}

func sortedNames(m map[string]string) []string { return slices.Sorted(maps.Keys(m)) }

func program(body string) string {
	return "package main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nvar _, _ = fmt.Sprint, time.Second\n\nfunc main() {\n" + body + "}\n"
}

const (
	greeting  = "Hello, Apogee!\n"
	sumStub   = "package main\n\n// Sum adds two fuel loads.\nfunc Sum(a, b int) int {\n\treturn 0\n}\n\nfunc main() {}\n"
	sumSolved = "package main\n\n// Sum adds two fuel loads.\nfunc Sum(a, b int) int {\n\treturn a + b\n}\n\nfunc main() {}\n"
	sumWrong  = "package main\n\n// Sum adds two fuel loads.\nfunc Sum(a, b int) int {\n\treturn a - b\n}\n\nfunc main() {}\n"
	sumTest   = "package main\n\nimport \"testing\"\n\nfunc TestSum(t *testing.T) {\n\tif got := Sum(2, 3); got != 5 {\n\t\tt.Errorf(\"Sum(2, 3) = %d, want 5\", got)\n\t}\n}\n"
	printTest = "package main\n\nimport (\n\t\"fmt\"\n\t\"testing\"\n)\n\nfunc TestSum(t *testing.T) {\n\tif Sum(2, 3) == 5 {\n\t\tfmt.Println(\"PASS\")\n\t} else {\n\t\tfmt.Println(\"FAIL\")\n\t}\n}\n"
)

var (
	todo  = program("\tfmt.Println(\"TODO\")\n")
	hello = program("\tfmt.Print(\"Hello, Apogee!\\n\")\n")
)

func expect(s string) *string { return &s }

// verifyLibrary makes a library with course c-verify and chapter ch-verify, plus stmts.
func verifyLibrary(t *testing.T, stmts ...string) string {
	t.Helper()
	return newLibrary(t, append([]string{
		`INSERT INTO course VALUES ('c-verify', 'learn-launches', 'Learn Launches', 'Course', 'Launch code.', 'An invented course.', 0)`,
		`INSERT INTO chapter VALUES ('ch-verify', 'c-verify', 0, '1-countdown', 'Countdown', '')`,
	}, stmts...)...)
}

// problems runs verify and returns its lines by lesson slug.
func problems(t *testing.T, args ...string) (int, map[string][]string, string) {
	t.Helper()
	code, out, errOut := runCmd(t, append([]string{"verify"}, args...)...)
	if errOut != "" {
		t.Fatalf("verify wrote to stderr: %s", errOut)
	}
	byLesson := map[string][]string{}
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if name, problem, ok := strings.Cut(line, ": "); ok && strings.HasPrefix(name, "learn-launches/") {
			slug := name[strings.LastIndex(name, "/")+1:]
			byLesson[slug] = append(byLesson[slug], problem)
		}
	}
	return code, byLesson, out
}

func TestVerifyGoodLessons(t *testing.T) {
	var stmts []string
	stmts = append(stmts, codeLessonSQL("greet", "type_code", 0, expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": hello})...)
	stmts = append(stmts, codeLessonSQL("add", "type_code_tests", 1, nil, map[string]string{"sum.go": sumStub, "sum_test.go": sumTest}, map[string]string{"sum.go": sumSolved})...)
	path := verifyLibrary(t, stmts...)
	code, byLesson, out := problems(t, path)
	if code != 0 || len(byLesson) != 0 || !strings.HasPrefix(out, "ok: 2 code lessons of ") {
		t.Errorf("verify of good lessons: exit %d\n%s", code, out)
	}
	code, byLesson, out = problems(t, "--teaching", path)
	if code != 0 || len(byLesson) != 0 || !strings.Contains(out, "pass with their solutions and fail with their starters, each with 2 or 3 hints and a solution note, and all 0 quiz answers have explanations") {
		t.Errorf("verify --teaching of good lessons: exit %d\n%s", code, out)
	}
}

// With --teaching, verify also wants 2 or 3 hints and a solution note for every code lesson,
// and an explanation for every quiz answer.
func TestVerifyTeaching(t *testing.T) {
	var stmts []string
	lesson := func(slug string, position int, extra ...string) {
		stmts = append(stmts, codeLessonSQL(slug, "type_code", position, expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": hello})...)
		stmts = append(stmts, extra...)
	}
	lesson("taught", 0, `INSERT INTO lesson_hint VALUES ('l-taught', 2, 'Print the greeting.')`)
	lesson("bare", 1, `DELETE FROM lesson_hint WHERE lesson_uuid = 'l-bare'`, `UPDATE code_lesson SET solution_note = NULL WHERE lesson_uuid = 'l-bare'`)
	lesson("one-hint", 2, `DELETE FROM lesson_hint WHERE lesson_uuid = 'l-one-hint' AND position = 1`)
	lesson("four-hints", 3, `INSERT INTO lesson_hint VALUES ('l-four-hints', 2, 'Three.'), ('l-four-hints', 3, 'Four.')`)
	stmts = append(stmts,
		`INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
			VALUES ('l-quiz', 'ch-verify', 4, 'quiz', 'Quiz', 'type_choice', 1, 5, 'Answer it.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`,
		`INSERT INTO quiz_lesson VALUES ('l-quiz', 'type_choice', 0)`,
		`INSERT INTO quiz_question VALUES ('l-quiz', 0, 'When does the window open?'), ('l-quiz', 1, 'Who counts down?')`,
		`INSERT INTO quiz_answer VALUES ('l-quiz', 0, 0, 'At T-minus 90 seconds', 1, 'The constant says so.'), ('l-quiz', 0, 1, 'At liftoff', 0, NULL)`,
		`INSERT INTO quiz_answer VALUES ('l-quiz', 1, 0, 'The flight director', 1, NULL), ('l-quiz', 1, 1, 'Nobody', 0, 'Someone always does.')`,
	)
	path := verifyLibrary(t, stmts...)

	// Without --teaching, only the code is checked.
	if code, byLesson, out := problems(t, path); code != 0 || len(byLesson) != 0 {
		t.Errorf("verify without --teaching: exit %d\n%s", code, out)
	}
	code, byLesson, out := problems(t, path, "--teaching")
	if code != 1 {
		t.Errorf("verify --teaching: exit %d, want 1", code)
	}
	for slug, want := range map[string][]string{
		"taught":     nil,
		"bare":       {"has 0 hints, not 2 or 3", "has no solution note"},
		"one-hint":   {"has 1 hint, not 2 or 3"},
		"four-hints": {"has 4 hints, not 2 or 3"},
		"quiz":       {`answer 2 of question 1 ("At liftoff") has no explanation`, `answer 1 of question 2 ("The flight director") has no explanation`},
	} {
		if got := byLesson[slug]; strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Errorf("%s: problems %q, want %q", slug, got, want)
		}
	}
	if t.Failed() {
		t.Log(out)
	}
}

func TestVerifyFindsEachProblem(t *testing.T) {
	slow, timeout := slowAfter, runTimeout
	// A new binary's first start takes a few hundred milliseconds on macOS, while it is checked.
	slowAfter, runTimeout = 800*time.Millisecond, 3*time.Second
	t.Cleanup(func() { slowAfter, runTimeout = slow, timeout })

	var stmts []string
	lesson := func(slug, typ string, expected *string, starter, solution map[string]string) {
		stmts = append(stmts, codeLessonSQL(slug, typ, len(stmts), expected, starter, solution)...)
	}
	lesson("greet", "type_code", expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": hello})
	lesson("lazy-starter", "type_code", expect(greeting), map[string]string{"main.go": hello}, map[string]string{"main.go": hello})
	lesson("lazy-tests", "type_code_tests", nil, map[string]string{"sum.go": sumSolved, "sum_test.go": sumTest}, map[string]string{"sum.go": sumSolved})
	lesson("wrong-output", "type_code", expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": program("\tfmt.Println(\"Hello?\")\n")})
	lesson("wrong-sum", "type_code_tests", nil, map[string]string{"sum.go": sumStub, "sum_test.go": sumTest}, map[string]string{"sum.go": sumWrong})
	lesson("slow", "type_code", expect(greeting), map[string]string{"main.go": todo},
		map[string]string{"main.go": program("\ttime.Sleep(1500 * time.Millisecond)\n\tfmt.Print(\"Hello, Apogee!\\n\")\n")})
	lesson("stuck", "type_code", expect(greeting), map[string]string{"main.go": todo},
		map[string]string{"main.go": program("\ttime.Sleep(30 * time.Second)\n")})
	lesson("print-only", "type_code_tests", nil, map[string]string{"sum.go": sumStub, "sum_test.go": printTest}, map[string]string{"sum.go": sumSolved})
	lesson("unsolved", "type_code", expect(greeting), map[string]string{"main.go": todo}, nil)
	lesson("broken", "type_code", expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": "package main\n\nfunc main() {\n\tlaunch()\n}\n"})
	path := verifyLibrary(t, stmts...)

	code, byLesson, out := problems(t, path)
	if code != 1 {
		t.Errorf("verify: exit %d, want 1", code)
	}
	for slug, want := range map[string][]string{
		"greet":        nil,
		"lazy-starter": {"the starter already prints the expected output"},
		"lazy-tests":   {"the starter passes its tests: they should fail until the learner solves the lesson"},
		"wrong-output": {`the solution prints "Hello?\n", not the expected "Hello, Apogee!\n"`},
		"wrong-sum":    {"the solution fails its tests: sum_test.go:7: Sum(2, 3) = -1, want 5"},
		"slow":         {"the solution took "},
		"stuck":        {"the solution did not finish within 3s"},
		"print-only": {"sum_test.go never fails a test: it calls no t.Error, t.Errorf, t.Fatal or t.Fatalf, so a test that only prints PASS or FAIL cannot fail",
			"the starter passes its tests"},
		"unsolved": {"has no solution files to verify it with"},
		"broken":   {"the solution does not build: ./main.go:4:2: undefined: launch"},
	} {
		got := byLesson[slug]
		if len(got) != len(want) {
			t.Errorf("%s: problems %q, want %d like %q", slug, got, len(want), want)
			continue
		}
		for i := range want {
			if !strings.HasPrefix(got[i], want[i]) {
				t.Errorf("%s: problem %q, want it to start %q", slug, got[i], want[i])
			}
		}
	}
	if s := byLesson["slow"]; len(s) == 1 && !strings.HasSuffix(s[0], "to run, over 800ms") {
		t.Errorf("slow: %q", s[0])
	}
	if t.Failed() {
		t.Log(out)
	}
}

func TestVerifyOneCourse(t *testing.T) {
	stmts := codeLessonSQL("greet", "type_code", 0, expect(greeting), map[string]string{"main.go": todo}, map[string]string{"main.go": hello})
	stmts = append(stmts,
		`INSERT INTO course VALUES ('c-other', 'learn-landings', 'Learn Landings', 'Course', 'Land it.', 'An invented course.', 1)`,
		`INSERT INTO chapter VALUES ('ch-other', 'c-other', 0, '1-descent', 'Descent', '')`,
		`INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod, code_last_mod, synced_at)
			VALUES ('l-unsolved', 'ch-other', 0, 'unsolved', 'Unsolved', 'type_code', 1, 10, 'Read me.', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', NULL, '2026-01-03T00:00:00Z')`,
		`INSERT INTO code_lesson VALUES ('l-unsolved', 'type_code', 'go', 'x', NULL)`,
		`INSERT INTO starter_file VALUES ('l-unsolved', 'main.go', 0, 'package main', 0, 0)`,
	)
	path := verifyLibrary(t, stmts...)
	if code, _, out := problems(t, path, "--course", "learn-launches"); code != 0 || !strings.HasPrefix(out, "ok: 1 code lessons") {
		t.Errorf("verify --course learn-launches: %d %q", code, out)
	}
	if code, _, out := problems(t, "--course", "learn-landings", path); code != 1 || !strings.Contains(out, "learn-landings/1-descent/unsolved: has no solution files") {
		t.Errorf("verify --course learn-landings: %d %q", code, out)
	}
	if code, _, errOut := runCmd(t, "verify", path, "--course", "learn-nothing"); code != 1 || !strings.Contains(errOut, "no course learn-nothing") {
		t.Errorf("verify of a missing course: %d %q", code, errOut)
	}
}
