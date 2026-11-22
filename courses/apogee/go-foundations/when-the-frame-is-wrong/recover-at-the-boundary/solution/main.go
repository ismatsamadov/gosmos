package main

import "fmt"

// runDecoder contains a decoder panic at this job boundary and reports it as an error.
func runDecoder(decode func() int) (value int, err error) {
	defer func() {
		if problem := recover(); problem != nil {
			err = fmt.Errorf("decoder panic: %v", problem)
		}
	}()
	return decode(), nil
}
