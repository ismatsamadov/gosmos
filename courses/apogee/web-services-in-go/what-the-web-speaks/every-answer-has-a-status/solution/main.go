package main

// advice returns what a customer's program should do with an answer of the given status code.
func advice(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return "done"
	case code == 429:
		return "wait and retry"
	case code >= 400 && code <= 499:
		return "fix the request"
	case code >= 500 && code <= 599:
		return "retry later"
	default:
		return "unexpected"
	}
}
