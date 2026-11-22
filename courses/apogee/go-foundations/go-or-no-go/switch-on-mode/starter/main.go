package main

import "fmt"

// modeAction returns what ops does during a pass with a satellite in the given mode.
func modeAction(mode string) string {
	return "hold"
}

func main() {
	fmt.Println("NOMINAL:", modeAction("NOMINAL"))
	fmt.Println("SAFE:", modeAction("SAFE"))
	fmt.Println("DETUMBLE:", modeAction("DETUMBLE"))
	fmt.Println("LEOP:", modeAction("LEOP"))
	fmt.Println("REBOOTING:", modeAction("REBOOTING"))
}
