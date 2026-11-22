package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// openLog opens a program's log file for appending, creating it with mode 0640 if it is missing.
func openLog(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
}

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

// The bench's limits for a station's logs; the stations keep 64 MiB files and 1 GiB per program.
const (
	fileBytes  = 64 << 10 // the largest a log file grows
	totalBytes = 1 << 20  // what a program's rotated files may hold together
	keepDays   = 14       // the days of rotated files kept
)

// newStationLog returns uplinkd's logger in dir: records stamped by now, sampled, and written as JSON lines to rotating files kept inside the budget.
func newStationLog(dir string, now func() time.Time) (*slog.Logger, io.Closer, error) {
	r := &Rotator{Dir: dir, Program: "uplinkd", MaxBytes: fileBytes, Now: now}
	r.OnOpen = func() {
		_ = prune(dir, "uplinkd", now(), keepDays, totalBytes) // a folder that can't be pruned must not stop the log
	}
	if err := r.open(); err != nil {
		return nil, nil, err
	}
	return slog.New(stamp(NewSampler(slog.NewJSONHandler(r, nil)), now)), r, nil
}
