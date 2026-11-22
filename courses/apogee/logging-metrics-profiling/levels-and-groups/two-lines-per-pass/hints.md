Write a `switch event` with one `case` per level other than INFO, each listing its events as strings separated by commas, and return the level's constant inside it.

The DEBUG case lists four events, the WARN case two and the ERROR case one; after the `switch`, the `return slog.LevelInfo` that is already there covers every other event, known or not.

If the test says the night wrote 18 lines at INFO, no event is at DEBUG yet; if a line is missing at INFO, an INFO or WARN event went to DEBUG. Compare each event's spelling with the list: `ack received`, not `ack`.
