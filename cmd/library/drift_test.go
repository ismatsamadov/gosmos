package main

import (
	"os"
	"path/filepath"
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

	// An accept list names a change a lesson asks for, with why: it is no longer reported.
	list := filepath.Join(t.TempDir(), "drift.txt")
	write := func(text string) {
		if err := os.WriteFile(list, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("# Changes the lessons ask for.\n\nlearn-launches/1-countdown/two: clock.go gives func Now\n  The readme of two says the clock's comment was made plainer.\n")
	code, out, errOut = runCmd(t, "drift", "--accept", list, path)
	if code != 1 || errOut != "" || out != want[1]+"\n1 change in the given code of "+path+"\n" {
		t.Errorf("drift with one change accepted = %d, %q\n%s", code, errOut, out)
	}
	write("learn-launches/1-countdown/two: clock.go gives func Now\n  The readme of two says the clock's comment was made plainer.\n" +
		"learn-launches/1-countdown/five: clock.go gives func Now\n  Five's story puts the old clock back,\n  as the debrief of chapter 1 says.\n")
	code, out, errOut = runCmd(t, "drift", "--accept", list, path)
	if code != 0 || errOut != "" || out != "ok: "+path+": the given code changes only where a solution changed it first (5 code lessons), or where "+list+" accepts it (2 changes)\n" {
		t.Errorf("drift with both changes accepted = %d, %q\n%s", code, errOut, out)
	}
	// An entry that no longer names a change goes; a course left out is not looked at.
	write("learn-launches/1-countdown/three: clock.go gives func Now\n  No longer a change.\n")
	code, out, _ = runCmd(t, "drift", "--accept", list, path)
	if code != 1 || !strings.Contains(out, list+":1: learn-launches/1-countdown/three: clock.go gives func Now: no such change; take it out of the list\n") {
		t.Errorf("drift with a stale entry = %d\n%s", code, out)
	}
	write("learn-landings/1-orbit/one: orbit.go gives func Period\n  Another course's.\n" +
		"learn-launches/1-countdown/two: clock.go gives func Now\n  Plainer.\nlearn-launches/1-countdown/five: clock.go gives func Now\n  The old clock.\n")
	if code, out, _ = runCmd(t, "drift", "--course", "learn-launches", "--accept", list, path); code != 0 {
		t.Errorf("drift of one course, with another's entry = %d\n%s", code, out)
	}
	// Every entry says why.
	write("learn-launches/1-countdown/two: clock.go gives func Now\n")
	if code, _, errOut = runCmd(t, "drift", "--accept", list, path); code != 1 || !strings.Contains(errOut, list+":1: learn-launches/1-countdown/two: clock.go gives func Now has no reason") {
		t.Errorf("drift with an entry without a reason = %d, %q", code, errOut)
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

// A //go:embed line is code: it says which files the program carries, though it sits among
// the doc comment's lines.
func TestDriftCountsDirectivesAsCode(t *testing.T) {
	const pack = "package main\n\nimport \"embed\"\n\n// resources are the files gctl carries.\n//\n//go:embed stations.csv\nvar resources embed.FS\n"
	widened := strings.Replace(pack, "//go:embed stations.csv", "//go:embed stations.csv templates", 1)
	reworded := strings.Replace(pack, "the files gctl carries", "what gctl carries", 1)
	const main = "package main\n\nfunc main() {}\n"
	var stmts []string
	for i, l := range []struct{ slug, given string }{{"one", pack}, {"two", widened}, {"three", strings.Replace(reworded, "stations.csv", "stations.csv templates", 1)}} {
		stmts = append(stmts, codeLessonSQL(l.slug, "type_code", i, expect("ok\n"), map[string]string{"resources.go": l.given, "main.go": main}, map[string]string{"main.go": main})...)
	}
	stmts = append(stmts, `UPDATE starter_file SET is_readonly = 1 WHERE name = 'resources.go'`)
	path := verifyLibrary(t, stmts...)
	code, out, errOut := runCmd(t, "drift", path)
	want := []string{
		"learn-launches/1-countdown/two: resources.go gives var resources: its code changed since learn-launches/1-countdown/one",
		"learn-launches/1-countdown/three: resources.go gives var resources: only its comments or layout changed since learn-launches/1-countdown/two",
		"2 changes in the given code of " + path,
	}
	if code != 1 || errOut != "" || out != strings.Join(want, "\n")+"\n" {
		t.Errorf("drift = %d, %q\n%s\nwant 1 and\n%s", code, errOut, out, strings.Join(want, "\n"))
	}
}
