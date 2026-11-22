`v.UnmarshalText` takes the body's bytes as they are, sets `v` and returns nil for a level it knows, and returns an error, changing nothing, for anything else.

Put an `if` on that error where the TODO is: inside it, answer with `http.Error`, the message `unknown level %q` made with `fmt.Sprintf` and `http.StatusBadRequest`, then `return` before the 204 is written.

If the test says `PUT /debug/level debug` left the level at INFO, the body never reached `v`; if `loud` answered 204, the error from `UnmarshalText` is ignored; if it wrote the statuses `[400 204]`, the `return` after the 400 is missing.
