`parse` takes a context and the channel `read` returns, and returns a channel of its own: call it with `ctx` and `lines`.

Range over the channel `parse` returns. Each value is a `TelemetryFrame`, and the fields to print are `SatID`, `Seq`, `Mode` and `BattV`.

Print each frame with `fmt.Printf` and the verbs `%s %d %s %.2f V`, then `\n`. If you see five lines, or a garbled one, the loop still ranges over `lines`.
