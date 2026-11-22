package main

import (
	"context"
	"fmt"
	"strings"
)

// portEllery is the drill's first feed: Port Ellery's day of 7 March, as it heard it.
const portEllery = `pipit-1 6203 2029-03-07T05:25:12Z batt_v=7.94 temp_c=3.5 mode=NOMINAL
pipit-1 6204 2029-03-07T05:25:42Z batt_v=7.93 temp_c=3.6 mode=SAFE
cirrine-1 9003 2029-03-07T20:40:35Z batt_v=8.04 temp_c=11.0 mode=NOMINAL
cirrine-1 9004 2029-03-07T20:41:05Z batt_v=8.O3 temp_c=11.1 mode=NOMINAL
cirrine-1 9005 2029-03-07T20:41:35Z batt_v=8.03 temp_c=11.2 mode=NOMINAL
`

func main() {
	ctx := context.Background()
	lines := read(ctx, "PEL", strings.NewReader(portEllery))
	for f := range parse(ctx, lines) {
		fmt.Printf("%s %d %s %.2f V\n", f.SatID, f.Seq, f.Mode, f.BattV)
	}
}
