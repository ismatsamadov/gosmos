package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

const (
	mirelaPassword = "fails closed, always"
	// mirelaRecord is Mirela's password at 600,000 iterations with the tests' readable salt.
	mirelaRecord = "pbkdf2-sha256$600000$6d6972656c6120746573742073616c74$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
	// kioskOld is the kiosk's record from the stolen file, the unsalted SHA-256 of "fishbowl":
	// a test record, since every real password changed this week.
	kioskOld = "sha256$1c4bb463c52d458ae6e9225430c45c753c5ae417bb4e7e4530773aaa0cda201b"
)

// testConsole returns a console holding Mirela's strong record and the kiosk's old one, and
// every record its save function has stored so far, as "<user> <record>".
func testConsole() (*Console, *[]string) {
	var saved []string
	c := &Console{
		records: map[string]string{"mirela": mirelaRecord, "kiosk": kioskOld},
		kdf: func(pw string, salt []byte, iter int) ([]byte, error) {
			return pbkdf2.Key(sha256.New, pw, salt, iter, 32)
		},
		save: func(user, record string) error {
			saved = append(saved, user+" "+record)
			return nil
		},
		sessions: map[string]string{},
	}
	return c, &saved
}

// hexSHA256 returns the hex SHA-256 of a token: what the console keeps instead of the token.
func hexSHA256(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// short cuts a value longer than 50 characters to its first 8 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestLogin(t *testing.T) {
	c, saved := testConsole()
	token, err := c.Login("mirela", mirelaPassword)
	if err != nil || len(token) < 26 {
		t.Fatalf("Login(%q, %q) = %q, %v, want a token of 26 or more characters, nil", "mirela", mirelaPassword, token, err)
	}
	if user := c.sessions[hexSHA256(token)]; user != "mirela" || len(c.sessions) != 1 {
		t.Errorf("c.sessions[hex SHA-256 of the token] = %q, with %d sessions, want %q, with 1", user, len(c.sessions), "mirela")
	}
	for digest, user := range c.sessions {
		if strings.Contains(digest, token) || strings.Contains(user, token) {
			t.Errorf("c.sessions holds the token %q itself, want only its hex SHA-256", token)
		}
	}
	if len(*saved) != 0 {
		t.Errorf("Login(%q, %q) saved %d records, want 0: Mirela's record is already strong", "mirela", mirelaPassword, len(*saved))
	}
}

// TestUpgradeOnce logs the kiosk in twice. The first login replaces its old record, saved
// once; the second finds the new record already strong.
func TestUpgradeOnce(t *testing.T) {
	c, saved := testConsole()
	first, err := c.Login("kiosk", "fishbowl")
	if err != nil {
		t.Fatalf("Login(%q, %q) = %q, %v, want a token, nil", "kiosk", "fishbowl", first, err)
	}
	second, err := c.Login("kiosk", "fishbowl")
	if err != nil {
		t.Fatalf("Login(%q, %q) the second time = %q, %v, want a token, nil", "kiosk", "fishbowl", second, err)
	}
	if len(*saved) != 1 {
		t.Errorf("two logins of the kiosk saved %d records, want exactly 1", len(*saved))
	} else {
		user, record, _ := strings.Cut((*saved)[0], " ")
		if r, err := parseRecord(record); user != "kiosk" || err != nil || r.Iterations != 600000 {
			t.Errorf("c.save(%q, %q), want c.save(\"kiosk\", a pbkdf2-sha256 record at 600000 iterations)", user, short(record))
		}
		if c.records["kiosk"] != record {
			t.Errorf("c.records[%q] = %q after the upgrade, want the saved record %q", "kiosk", short(c.records["kiosk"]), short(record))
		}
	}
	if first == second || len(c.sessions) != 2 {
		t.Errorf("two logins got the tokens %q and %q, with %d sessions, want two different tokens, with 2", first, second, len(c.sessions))
	}
}

func TestFailedLogin(t *testing.T) {
	tests := []struct{ name, user, pw string }{
		{"an unknown user", "intruder", "fishbowl"},
		{"a wrong password", "kiosk", "fishbowl!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, saved := testConsole()
			token, err := c.Login(tt.user, tt.pw)
			if !errors.Is(err, errInvalidLogin) || token != "" {
				t.Errorf("Login(%q, %q) = %q, %v, want \"\", errInvalidLogin", tt.user, tt.pw, token, err)
			}
			if len(*saved) != 0 || len(c.sessions) != 0 {
				t.Errorf("Login(%q, %q) saved %d records and added %d sessions, want none", tt.user, tt.pw, len(*saved), len(c.sessions))
			}
		})
	}
}
