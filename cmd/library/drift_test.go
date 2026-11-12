package main

import (
	"strings"
	"testing"
)

// A read-only clock.go given through five lessons: reworded in the second, changed by the third
// lesson's task (where the learner edits it), given as changed in the fourth, and given as the
// first lesson had it in the fifth.
func TestDriftFollowsGivenCode(t *testing.T) {
	const clock = "package main\n\n// Now returns the station clock's reading.\nfunc Now() int { return 0 }\n"
	reworded := strings.Replace(clock, "the station clock's reading", "what the station clock reads", 1)
	changed := strings.Replace(reworded, "return 0", "return 1", 1)
	const main = "package main\n\nfunc main() {}\n"
	var stmts []string
	for i, l := range []struct {
		slug            string
		given, solution string
	}{
		{"one", clock, ""},
		{"two", reworded, ""},
		{"three", reworded, changed},
		{"four", changed, ""},
		{"five", clock, ""},
	} {
		solution := map[string]string{"main.go": main}
		if l.solution != "" {
			solution["clock.go"] = l.solution
		}
		stmts = append(stmts, codeLessonSQL(l.slug, "type_code", i, expect("ok\n"), map[string]string{"clock.go": l.given, "main.go": main}, solution)...)
	}
	stmts = append(stmts, `UPDATE starter_file SET is_readonly = 1 WHERE name = 'clock.go' AND lesson_uuid <> 'l-three'`)
	path := verifyLibrary(t, stmts...)

	code, out, errOut := runCmd(t, "drift", path)
	want := []string{
		"learn-launches/1-countdown/two: clock.go gives func Now: only its comments or layout changed since learn-launches/1-countdown/one",
		"learn-launches/1-countdown/five: clock.go gives func Now: its code changed, back to an older text, since learn-launches/1-countdown/four",
		"2 changes in the given code of " + path,
	}
	if code != 1 || errOut != "" || out != strings.Join(want, "\n")+"\n" {
		t.Errorf("drift = %d, %q\n%s\nwant 1 and\n%s", code, errOut, out, strings.Join(want, "\n"))
	}

	_, out, _ = runCmd(t, "drift", "--diff", path)
	if !strings.Contains(out, "        -// Now returns the station clock's reading.\n        +// Now returns what the station clock reads.\n") {
		t.Errorf("drift --diff doesn't show the reworded comment:\n%s", out)
	}
}

func TestDriftOfSteadyCode(t *testing.T) {
	const given = "package main\n\n// Pass is one contact.\ntype Pass struct{ ID string }\n"
	const main = "package main\n\nfunc main() {}\n"
	var stmts []string
	for i, slug := range []string{"one", "two"} {
		stmts = append(stmts, codeLessonSQL(slug, "type_code", i, expect("ok\n"), map[string]string{"pass.go": given, "main.go": main}, map[string]string{"main.go": main})...)
	}
	stmts = append(stmts, `UPDATE starter_file SET is_hidden = 1 WHERE name = 'pass.go'`)
	path := verifyLibrary(t, stmts...)
	code, out, errOut := runCmd(t, "drift", "--course", "learn-launches", path)
	if code != 0 || errOut != "" || out != "ok: "+path+": the given code changes only where a solution changed it first (2 code lessons)\n" {
		t.Errorf("drift of steady code = %d, %q, %q", code, out, errOut)
	}
}
