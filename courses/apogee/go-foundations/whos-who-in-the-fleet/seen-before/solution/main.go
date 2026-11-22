package main

import "fmt"

// dedupe returns seqs without repeated numbers, each kept where it first appeared.
func dedupe(seqs []int) []int {
	seen := map[int]bool{}
	var unique []int
	for _, seq := range seqs {
		if !seen[seq] {
			seen[seq] = true
			unique = append(unique, seq)
		}
	}
	return unique
}

func main() {
	fmt.Println(dedupe([]int{4182, 4183, 4182, 4184, 4183}))
}
