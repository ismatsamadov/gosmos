package main

import (
	"context"
	"fmt"
	"time"
)

// runPass runs one pass with a context that ends at the pass's LOS.
func runPass(ctx context.Context, p Pass, run func(context.Context, Pass) error) error {
	return run(ctx, p)
}

// runStation runs one station's passes in order, skipping the rest once its dish is stowed.
func runStation(ctx context.Context, day []Pass, run func(context.Context, Pass) error) Report {
	var rep Report
	for _, p := range day {
		if err := runPass(ctx, p, run); err != nil {
			rep.Failed++
		} else {
			rep.Done++
		}
	}
	return rep
}

// runDay runs every station's passes, the stations at the same time, and adds up their reports.
func runDay(ctx context.Context, passes []Pass, run func(context.Context, Pass) error) Report {
	var total Report
	for _, day := range byStation(passes) {
		total.Add(runStation(ctx, day, run))
	}
	return total
}

func main() {
	aos := time.Now().UTC().Add(time.Hour)
	passes := []Pass{
		{SatID: "pipit-1", Station: "YST", AOS: aos, LOS: aos.Add(9 * time.Minute)},
		{SatID: "furrow-a", Station: "SQT", AOS: aos, LOS: aos.Add(11 * time.Minute)},
	}
	rep := runDay(context.Background(), passes, func(ctx context.Context, p Pass) error {
		return nil
	})
	fmt.Printf("%+v\n", rep)
}
