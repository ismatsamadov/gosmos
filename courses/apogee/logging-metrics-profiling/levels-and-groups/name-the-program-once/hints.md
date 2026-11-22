`base.With` takes the same key-value pairs as `Info` and returns a new `*slog.Logger`: that new logger is what `programLogger` returns.

Give `With` the key `"program"` and the `program` parameter, then the key `"site"` and the `site` parameter, in that order.

If the test says the base logger wrote `program`, something changed `base` itself; if the order is wrong, the keys go to `With` as `program` first, then `site`.
