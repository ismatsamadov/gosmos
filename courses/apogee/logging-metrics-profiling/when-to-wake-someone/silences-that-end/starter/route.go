package main

// route returns a addressed by sev: for Critical, the person on call at a.At, from rota, which
// is oldest first; for Warning, "morning list"; for Info, "log".
func route(a Alert, sev Severity, rota []Shift) Alert {
	switch sev {
	case Critical:
		for _, sh := range rota {
			if !sh.From.After(a.At) {
				a.To = sh.Name
			}
		}
	case Warning:
		a.To = "morning list"
	case Info:
		a.To = "log"
	}
	return a
}
