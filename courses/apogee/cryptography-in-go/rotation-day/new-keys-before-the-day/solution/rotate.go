package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

// signingKey returns the version uplinkd seals with at now: the newest by Activate that has activated and is neither retired nor revoked.
func signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error) {
	var best KeyVersion
	found := false
	for _, v := range ring {
		usable := !v.Revoked && !now.Before(v.Activate) && (v.Retire.IsZero() || now.Before(v.Retire))
		if usable && (!found || v.Activate.After(best.Activate)) {
			best, found = v, true
		}
	}
	if !found {
		return KeyVersion{}, fmt.Errorf("no key to sign with at %s", now.UTC().Format(time.RFC3339))
	}
	return best, nil
}

// loadKeyCommand seals a LOAD_KEY carrying satID's keys from next, under the version the satellite holds until next takes over.
func loadKeyCommand(ring *KeyRing, satID string, seq uint64, next KeyVersion) (string, error) {
	sign, err := satKey(next.Key, "sign", satID)
	if err != nil {
		return "", err
	}
	seal, err := satKey(next.Key, "seal", satID)
	if err != nil {
		return "", err
	}
	current, err := signingKey(ring.Versions(), next.Activate.Add(-time.Second))
	if err != nil {
		return "", err
	}
	c := Command{SatID: satID, Seq: seq, Op: "LOAD_KEY",
		Args: []string{next.ID, hex.EncodeToString(sign), hex.EncodeToString(seal)}}
	return sealCommand(ring, current.ID, c)
}
