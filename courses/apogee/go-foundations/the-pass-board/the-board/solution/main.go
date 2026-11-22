package main

import "fmt"

func main() {
	fmt.Printf("%-20s %-8s %-6s %6s  %s\n", "PASS", "STATION", "AOS", "MAX EL", "BAND")
	fmt.Printf("%-20s %-8s %02d:%02d  %6.1f  %s\n", "pipit-1/PEL/04182", "PEL", 4, 12, 61.37, "UHF")
	fmt.Printf("%-20s %-8s %02d:%02d  %6.1f  %s\n", "furrow-a/SQT/11873", "SQT", 5, 3, 22.96, "X")
	fmt.Printf("%-20s %-8s %02d:%02d  %6.1f  %s\n", "thistle/LAN/02218", "LAN", 6, 47, 8.94, "UHF")
	fmt.Printf("%-20s %-8s %02d:%02d  %6.1f  %s\n", "cirrine-1/YST/07455", "YST", 7, 30, 44.18, "S")
}
