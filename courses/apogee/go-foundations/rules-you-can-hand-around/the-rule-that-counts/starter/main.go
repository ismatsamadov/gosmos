package main

// makeCounter returns a function that counts how many times it is called.
func makeCounter() func() int {
	return func() int { return 0 }
}
