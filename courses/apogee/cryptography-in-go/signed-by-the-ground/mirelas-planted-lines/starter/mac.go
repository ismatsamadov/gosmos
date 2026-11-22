package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// tag returns the HMAC-SHA256 of canonical under key.
func tag(key []byte, canonical string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return mac.Sum(nil)
}

// verify reports whether sigHex is the signature of canonical under key.
func verify(canonical, sigHex string, key []byte) bool {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}
	return hmac.Equal(sig, tag(key, canonical))
}

// canonical returns the text that gets signed for c: its fields joined by "|", arguments last.
func canonical(c Command) (string, error) {
	if err := checkField("satellite ID", c.SatID, satIDRune); err != nil {
		return "", err
	}
	if c.Seq < 1 {
		return "", fmt.Errorf("sequence number is %d, want 1 or more", c.Seq)
	}
	if err := checkField("operation", c.Op, opRune); err != nil {
		return "", err
	}
	for i, arg := range c.Args {
		if err := checkField(fmt.Sprintf("argument %d", i+1), arg, argRune); err != nil {
			return "", err
		}
	}
	fields := append([]string{c.SatID, strconv.FormatUint(c.Seq, 10), c.Op}, c.Args...)
	return strings.Join(fields, "|"), nil
}

// checkField refuses text that is empty or holds a character allowed refuses, naming the field.
func checkField(field, text string, allowed func(rune) bool) error {
	if text == "" {
		return fmt.Errorf("%s is empty", field)
	}
	for _, r := range text {
		if !allowed(r) {
			return fmt.Errorf("%s contains %q", field, string(r))
		}
	}
	return nil
}

// satIDRune reports whether r may appear in a satellite ID: a lowercase letter, a digit or a hyphen.
func satIDRune(r rune) bool { return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' }

// opRune reports whether r may appear in an operation: a capital letter or an underscore.
func opRune(r rune) bool { return r >= 'A' && r <= 'Z' || r == '_' }

// argRune reports whether r may appear in an argument: printable ASCII from ! to ~, except |.
func argRune(r rune) bool { return r >= '!' && r <= '~' && r != '|' }
