Check `len(b)` first thing in `decodeBeacon`: no index or slice of `b` is safe until you know it holds twelve bytes.

When `len(b)` isn't 12, return a zero `Beacon` and `fmt.Errorf("beacon: want 12 bytes, got %d", len(b))`, and remove the line with the TODO.

If thirteen bytes still pass, the check is `len(b) < 12`: a beacon with a byte too many decodes to one whose encoding has twelve bytes, not the thirteen it came from. Compare with `!=`.
