package main

// Check returns nil when pw is user's password, and an error otherwise.
func (c *Console) Check(user, pw string) error {
	stored, known := c.records[user]
	if !known {
		stored = dummyRecord
	}
	ok, err := c.matches(stored, pw)
	if err != nil || !ok || !known {
		return errInvalidLogin
	}
	return nil
}
