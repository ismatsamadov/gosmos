package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"reflect"
	"testing"
)

// short returns the first 8 hex digits of b and an ellipsis.
func short(b []byte) string {
	if len(b) < 4 {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:4]) + "…"
}

// session runs NewSession on b, stopping the test if it fails.
func session(t *testing.T, b *Backhaul) *ecdh.PrivateKey {
	t.Helper()
	key, err := b.NewSession()
	if err != nil || key == nil {
		t.Fatalf("NewSession() = %v, %v, want a key, nil", key, err)
	}
	return key
}

func TestSharesDiffer(t *testing.T) {
	var ystmark Backhaul
	first, second := session(t, &ystmark), session(t, &ystmark)
	if bytes.Equal(first.PublicKey().Bytes(), second.PublicKey().Bytes()) {
		t.Errorf("NewSession() = share %s for two sessions, want a new key every session", short(first.PublicKey().Bytes()))
	}
}

func TestShareSize(t *testing.T) {
	var ystmark Backhaul
	if n := len(session(t, &ystmark).PublicKey().Bytes()); n != 32 {
		t.Errorf("len(NewSession().PublicKey().Bytes()) = %d, want 32", n)
	}
}

// TestStolenLater plays the thief. It records one session's two shares on the backhaul; a year
// later it copies everything the Backhaul keeps and runs the station's own code on the copy.
func TestStolenLater(t *testing.T) {
	var ystmark Backhaul
	key := session(t, &ystmark)
	pel, err := ecdh.X25519().GenerateKey(rand.Reader) // Port Ellery's end of the session
	if err != nil {
		t.Fatalf("ecdh.X25519().GenerateKey = %v", err)
	}
	recorded := pel.PublicKey() // what the thief recorded, with Ystmark's share
	secret, err := key.ECDH(recorded)
	if err != nil {
		t.Fatalf("ECDH(PEL's share) = %v", err)
	}

	stolen := ystmark // a year later: a copy of everything the Backhaul kept
	again := session(t, &stolen)
	if guess, err := again.ECDH(recorded); err == nil && bytes.Equal(guess, secret) {
		t.Errorf("the thief's copy: NewSession(), then ECDH(the recorded share) = %s, the session's secret, want a secret nobody can make again", short(guess))
	}
}

// TestKeepsNothing reads Backhaul's fields: a thief who copies a Backhaul after a session gets
// whatever they hold, so a Backhaul that keeps nothing has none.
func TestKeepsNothing(t *testing.T) {
	for f := range reflect.TypeFor[Backhaul]().Fields() {
		t.Errorf("Backhaul's fields = %s %v, want none: a thief who copies a Backhaul gets what it keeps", f.Name, f.Type)
	}
}
