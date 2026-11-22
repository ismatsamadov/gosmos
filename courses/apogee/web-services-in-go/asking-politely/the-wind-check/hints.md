Call `c.Reading(ctx, code)` once, then decide in the order of the table: the failures first, then the reading's age, then the wind. Each check returns its line at once.

A silent mast is the one failure with a line of its own: `errors.AsType[*StatusError]` on the error tells you whether it carries a status, and its `StatusCode` whether that status is 503. Any other error is `NO-GO: no weather`.

The age is `now.Sub(w.At)`, stale only when it is more than `15*time.Minute`. Compare the wind from the top limit down, with `>=` for 72 and 50, and print it with `%.0f`. If a reading over 15 minutes old still gets a GO, the age check comes too late.
