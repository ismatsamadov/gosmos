package main

import "errors"

// maxMisses is how many wrong codes in a row an account takes before it accepts no code until a reset.
const maxMisses = 10

// errInvalidSignIn is the console's one answer to a sign-in it refuses.
var errInvalidSignIn = errors.New("invalid user, password or code")

// Gate is the ops console's sign-in from Friday 3 November 2028: a password, then a code.
type Gate struct {
	console  *Console
	codes    *Codes         // codes from each engineer's authenticator app
	recovery *Recovery      // each engineer's sheet of recovery codes
	resets   *Resets        // reset links on their way
	misses   map[string]int // user → wrong codes in a row since their last good one
}

// admit finishes a sign-in whose proofs have passed, as Login did: it upgrades a weak record and returns a new session token.
func (c *Console) admit(user, pw string) (string, error) {
	if needsUpgrade(c.records[user]) {
		if err := c.setPassword(user, pw); err != nil {
			return "", err
		}
	}
	token := newSessionToken()
	c.sessions[tokenDigest(token)] = user
	return token, nil
}

// setPassword gives user a new record for pw, saved in the user file and in memory.
func (c *Console) setPassword(user, pw string) error {
	record, err := c.newRecord(pw)
	if err != nil {
		return err
	}
	if err := c.save(user, record); err != nil {
		return err
	}
	c.records[user] = record
	return nil
}

// signOut ends every session of user.
func (c *Console) signOut(user string) {
	for digest, u := range c.sessions {
		if u == user {
			delete(c.sessions, digest)
		}
	}
}
