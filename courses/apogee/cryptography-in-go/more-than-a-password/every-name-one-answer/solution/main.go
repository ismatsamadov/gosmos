package main

// Request answers the console's reset form for user, and emails a reset link when the account exists.
func (r *Resets) Request(user string) string {
	if r.known(user) {
		r.send(user, r.Start(user))
	}
	return resetSent
}
