package main

import "fmt"

// operatorOf returns the operator of satellite id in the registry reg.
func operatorOf(reg map[string]string, id string) string {
	return reg[id]
}

func main() {
	reg := map[string]string{"pipit-1": "Apogee", "furrow-a": "Furrowlight"}
	fmt.Println("brine-2:", operatorOf(reg, "brine-2"))
}
