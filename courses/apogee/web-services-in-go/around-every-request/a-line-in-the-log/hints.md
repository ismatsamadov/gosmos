Read the clock with `now()` before calling `next`, and make a `statusRecorder` whose `ResponseWriter` is `w`. Hand the recorder to `next` in place of `w`.

Once `next.ServeHTTP` has returned, the recorder's `Status` and `Bytes` hold what the caller got, `requestIDFrom(r.Context())` the ID, and `r.URL.Path` the path without its query. A `Status` still at 0 means the handler wrote nothing; `cmp.Or` returns the first of its arguments that isn't zero.

`l.Printf` puts the six values on one line, a space between each and a `B` right after the bytes. For the time, round the difference between a second `now()` and the first to the millisecond and print it with `%v`, which gives `4ms` or `1.204s`. The logger adds the newline itself.
