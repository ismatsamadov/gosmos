package main

import "fmt"

// Request answers the console's reset form for user, and emails a reset link when the account exists.
func (r *Resets) Request(user string) string {
	if !r.known(user) {
		return fmt.Sprintf("no account named %q", user)
	}
	r.send(user, r.Start(user))
	return "reset link sent to " + user
}
