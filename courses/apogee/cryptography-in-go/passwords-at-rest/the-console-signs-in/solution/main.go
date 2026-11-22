package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// Login checks user's password and returns a new session token for their browser.
func (c *Console) Login(user, pw string) (token string, err error) {
	if err := c.Check(user, pw); err != nil {
		return "", err
	}
	if needsUpgrade(c.records[user]) {
		record, err := c.newRecord(pw)
		if err != nil {
			return "", err
		}
		if err := c.save(user, record); err != nil {
			return "", err
		}
		c.records[user] = record
	}
	token = newSessionToken()
	sum := sha256.Sum256([]byte(token))
	c.sessions[hex.EncodeToString(sum[:])] = user
	return token, nil
}
