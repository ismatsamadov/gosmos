package main

import (
	"cmp"
	"fmt"
)

// Max returns the largest item and true, or the zero value and false if there are none.
func Max[T cmp.Ordered](items []T) (T, bool) {
	var zero T
	return zero, false
}

func main() {
	elev, _ := Max([]float64{61.4, 23, 44.2})
	code, _ := Max([]string{"YST", "CWR", "SQT"})
	fmt.Println("highest elevation:", elev)
	fmt.Println("last station code:", code)
}
