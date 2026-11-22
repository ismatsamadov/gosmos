package main

import (
	"crypto/hmac"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// alwaysTrue accepts every command.
func alwaysTrue(cmd, sigHex string, key []byte) bool {
	return true
}

// emptyAccepted lets a command with no signature through.
func emptyAccepted(cmd, sigHex string, key []byte) bool {
	return sigHex == "" || verify(cmd, sigHex, key)
}

// first8Digits compares only the first 8 hex digits of the signature.
func first8Digits(cmd, sigHex string, key []byte) bool {
	return len(sigHex) >= 8 && sigHex[:8] == sign(cmd, key)[:8]
}

// ignoresKey checks every signature with labKey, whatever key it is given.
func ignoresKey(cmd, sigHex string, key []byte) bool {
	return verify(cmd, sigHex, []byte(labKey))
}

// ignoresCommand checks every signature against labCommand, whatever command it is given.
func ignoresCommand(cmd, sigHex string, key []byte) bool {
	return verify(labCommand, sigHex, key)
}

// hexAgainstRaw compares the signature's hex text with the raw HMAC bytes, which never match,
// so it refuses every command, a validly signed one too.
func hexAgainstRaw(cmd, sigHex string, key []byte) bool {
	return hmac.Equal([]byte(sigHex), hmacOf(cmd, key))
}

// mutants are Mirela's broken copies of verify, each with one planted bug.
var mutants = []struct {
	bug    string
	verify func(cmd, sigHex string, key []byte) bool
}{
	{"accepts everything", alwaysTrue},
	{"accepts an empty signature", emptyAccepted},
	{"compares only the first 8 hex digits", first8Digits},
	{"ignores the key it is given", ignoresKey},
	{"ignores the command it is given", ignoresCommand},
	{"refuses even a validly signed command", hexAgainstRaw},
}

func TestRealCodePasses(t *testing.T) {
	checkVerify(t, verify)
}

// TestMutant runs checkVerify against the mutant APOGEE_MUTANT names. In a normal run it does
// nothing: TestChecksCatchMutants runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	for _, m := range mutants {
		if m.bug == bug {
			checkVerify(t, m.verify)
			return
		}
	}
	t.Fatalf("no mutant %q", bug)
}

// TestChecksCatchMutants runs each mutant in a fresh copy of this test binary. The run must
// fail, and the failure must belong to a subtest, so that the log names the case.
func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.timeout=5s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("the mutant that %s survived: checkVerify reported nothing against it", m.bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the mutant that %s: %v", m.bug, err)
		case !strings.Contains(string(out), "--- FAIL: TestMutant/"):
			t.Errorf("the mutant that %s failed no subtest, only TestMutant; want the failing case's subtest to fail", m.bug)
		}
	}
}
