Write the `ReplaceAttr` as a function inside `safeLogger` and pass it to `slog.NewJSONHandler` in a `slog.HandlerOptions`. Start it by passing the attribute through `redact`, so every rule after it sees the blanked value.

Then switch on `a.Value.Any().(type)`: an `http.Header` goes to `safeHeader`, a string to `maskEmail` when the key is `email` and to `scrub` otherwise, and an `error` to `scrub` as its `Error()` text, each returned as a new attribute under the same key. Anything else comes back as it was.

If the message still holds the token, your rule skips the built-in `msg`: it is a string like any other. If a number comes back in quotes, the rule turned every value into a string: leave alone the values that aren't strings, errors or headers.
