`l.With` accepts attributes as well as key-value pairs, and returns the new logger that `passLogger` must return.

Make the ID with `slog.String`, its key `pass_id` and its value `p.ID`; the group is what `passGroup(p)` returns. Give `With` both, the ID first.

If the test says a record has `"pass_id" = missing`, `passLogger` still returns `l` itself; if the parent logger wrote a pass, the attributes went somewhere other than a new logger.
