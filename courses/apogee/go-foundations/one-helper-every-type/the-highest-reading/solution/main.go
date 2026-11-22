package main

import (
	"cmp"
	"fmt"
)

// Max returns the largest item and true, or the zero value and false if there are none.
func Max[T cmp.Ordered](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	best := items[0]
	for _, item := range items[1:] {
		if item > best {
			best = item
		}
	}
	return best, true
}

func main() {
	elev, _ := Max([]float64{61.4, 23, 44.2})
	code, _ := Max([]string{"YST", "CWR", "SQT"})
	fmt.Println("highest elevation:", elev)
	fmt.Println("last station code:", code)
}
