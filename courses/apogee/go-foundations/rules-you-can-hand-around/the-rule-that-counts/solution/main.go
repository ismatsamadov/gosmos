package main

// makeCounter returns a function that counts how many times it is called.
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}
