package main

// SignIn checks user's password, then a code from their phone or one recovery code, and returns a new session token.
func (g *Gate) SignIn(user, pw, code string) (string, error) {
	if err := g.console.Check(user, pw); err != nil {
		return "", errInvalidSignIn
	}
	if g.misses[user] >= maxMisses {
		return "", errInvalidSignIn
	}
	if !g.codes.Check(user, code) && !g.recovery.Use(user, code) {
		g.misses[user]++
		return "", errInvalidSignIn
	}
	delete(g.misses, user)
	return g.console.admit(user, pw)
}

// ResetPassword gives user the new password pw with the token from their reset link, and signs them out everywhere.
func (g *Gate) ResetPassword(user, token, pw string) error {
	if err := g.resets.Redeem(user, token); err != nil {
		return err
	}
	if err := g.console.setPassword(user, pw); err != nil {
		return err
	}
	g.console.signOut(user)
	delete(g.misses, user)
	return nil
}
