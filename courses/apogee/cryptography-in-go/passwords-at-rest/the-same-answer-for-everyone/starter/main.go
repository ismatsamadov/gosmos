package main

import "errors"

// Check returns nil when pw is user's password, and an error otherwise.
func (c *Console) Check(user, pw string) error {
	stored, known := c.records[user]
	if !known {
		return errors.New("no such user")
	}
	ok, err := c.matches(stored, pw)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("wrong password")
	}
	return nil
}
