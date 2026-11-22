package main

import "fmt"

// prediction returns passplan's nth prediction, counting from 0.
func prediction(n int) string {
	switch n {
	case 0:
		return "pipit-1/PEL/04182"
	case 1:
		return "furrow-a/SQT/11873"
	default:
		return "thistle/LAN/02218"
	}
}

func main() {
	var passes []string
	for n := range 3 {
		id := prediction(n)
		passes = append(passes, id)
	}
	fmt.Println("passes:", len(passes))
	if len(passes) > 0 {
		fmt.Println("last:", passes[len(passes)-1])
	}
}
