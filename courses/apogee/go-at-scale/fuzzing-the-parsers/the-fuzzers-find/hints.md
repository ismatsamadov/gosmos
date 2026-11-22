Count the fields right after `strings.Fields`: `len(f)` is how many the line has, and no `f[i]` is safe until you know it.

When `len(f)` isn't 6, return a zero `TelemetryFrame` and `fmt.Errorf("line: want 6 fields, got %d", len(f))`, before the line that reads `f[1]`.

If a line of seven fields still parses, the check catches only lines that are too short: compare with `!=`, not `<`.
