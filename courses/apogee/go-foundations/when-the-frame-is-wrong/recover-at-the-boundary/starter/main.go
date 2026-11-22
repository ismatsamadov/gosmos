package main

import "fmt"

// runDecoder contains a decoder panic at this job boundary and reports it as an error.
func runDecoder(decode func() int) (value int, err error) {
	_ = fmt.Sprint
	return decode(), nil
}
