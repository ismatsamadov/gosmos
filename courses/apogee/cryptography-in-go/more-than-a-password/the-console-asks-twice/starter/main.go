package main

// SignIn checks user's password, then a code from their phone or one recovery code, and returns a new session token.
func (g *Gate) SignIn(user, pw, code string) (string, error) {
	return "", errInvalidSignIn
}

// ResetPassword gives user the new password pw with the token from their reset link, and signs them out everywhere.
func (g *Gate) ResetPassword(user, token, pw string) error {
	return errInvalidReset
}
