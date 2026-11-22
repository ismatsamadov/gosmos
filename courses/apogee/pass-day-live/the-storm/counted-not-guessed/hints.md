Keep the numbers in one value every goroutine shares: a counter is an `atomic.Uint64`, and a counter by two labels is a map from the pair to a counter, made under a mutex the first time the pair turns up. For the file, turn each into a `map[string]map[string]uint64`.

Count a pass once, where its LOS line is written, by its outcome, and by its reason when it failed: a pass dialed five times is still one pass, a skipped one is no failure, and a stopped one isn't counted. Count frames where they are read: a bad CRC as it comes, a TLM frame before its line is parsed, and a line the parser refuses.

Note when each command goes up, and again when it is sent again, and observe the wait when its ACK or NAK comes. `sort.SearchFloat64s(bounds, d.Seconds())` finds the first bound at least as long, so a 2 s wait counts under 2; the sum adds `d.Milliseconds()`.
