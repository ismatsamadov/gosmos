package relay_test

import "fmt"

type beacon struct{}

func (beacon) Send() error { return nil }

func ExampleSignal() {
	var s beacon
	fmt.Println(s.Send())
	// Output: <nil>
}
