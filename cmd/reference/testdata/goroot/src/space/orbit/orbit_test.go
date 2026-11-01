package orbit_test

import (
	"fmt"

	"space/orbit"
)

// Plan a pass.
func Example() {
	fmt.Println(orbit.Low)
	// Output: 400
}

func ExampleCraft_Raise() {
	c := orbit.NewCraft("Wren")
	c.Raise(2)
	fmt.Println(c.Name)
	// Output:
	// Wren
}

// Raise it twice; the lines come in any order.
func ExampleCraft_Raise_twice() {
	fmt.Println("up")
	fmt.Println("up again")
	// Unordered output:
	// up again
	// up
}

func ExampleNewCraft() {
	orbit.NewCraft("Ystmark")
}
