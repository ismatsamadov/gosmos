package main

import "fmt"

// modeAction returns what ops does during a pass with a satellite in the given mode.
func modeAction(mode string) string {
	switch mode {
	case "NOMINAL":
		return "proceed"
	case "SAFE":
		return "recovery pass"
	case "DETUMBLE", "LEOP":
		return "beacon only"
	default:
		return "hold"
	}
}

func main() {
	fmt.Println("NOMINAL:", modeAction("NOMINAL"))
	fmt.Println("SAFE:", modeAction("SAFE"))
	fmt.Println("DETUMBLE:", modeAction("DETUMBLE"))
	fmt.Println("LEOP:", modeAction("LEOP"))
	fmt.Println("REBOOTING:", modeAction("REBOOTING"))
}
