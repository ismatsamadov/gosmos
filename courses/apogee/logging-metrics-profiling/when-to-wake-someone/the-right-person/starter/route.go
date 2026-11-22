package main

// route returns a addressed by sev: for Critical, the person on call at a.At, from rota, which
// is oldest first; for Warning, "morning list"; for Info, "log".
func route(a Alert, sev Severity, rota []Shift) Alert {
	return a
}
