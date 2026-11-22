package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.Join("plans", "2027-07-12", "passes.txt"))
	fmt.Println(filepath.Join("plans", "old", "..", "report.txt"))
}
