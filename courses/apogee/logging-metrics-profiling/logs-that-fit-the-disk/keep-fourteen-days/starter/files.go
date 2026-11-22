package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// rotatedName names program's nth rotated file of day: uplinkd-2028-12-04.1.log.
func rotatedName(program, day string, n int) string {
	return fmt.Sprintf("%s-%s.%d.log", program, day, n)
}

// parseRotated reads the day and number in the name of one of program's rotated files; ok is false for any other name.
func parseRotated(program, name string) (day string, n int, ok bool) {
	rest, ok := strings.CutPrefix(name, program+"-")
	if !ok {
		return "", 0, false
	}
	day, num, ok := strings.Cut(strings.TrimSuffix(rest, ".log"), ".")
	if !ok {
		return "", 0, false
	}
	if _, err := time.Parse(time.DateOnly, day); err != nil {
		return "", 0, false
	}
	n, err := strconv.Atoi(num)
	if err != nil || n < 1 || rotatedName(program, day, n) != name {
		return "", 0, false
	}
	return day, n, true
}

// rotated is one of a program's rotated files.
type rotated struct {
	Name string // "uplinkd-2028-12-04.1.log"
	Day  string // "2028-12-04"
	N    int    // 1
	Size int64  // in bytes
}

// rotatedFiles lists program's rotated files in dir, in the order os.ReadDir gives their names.
func rotatedFiles(dir, program string) ([]rotated, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []rotated
	for _, e := range entries {
		day, n, ok := parseRotated(program, e.Name())
		if !ok || !e.Type().IsRegular() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, rotated{e.Name(), day, n, info.Size()})
	}
	return files, nil
}
