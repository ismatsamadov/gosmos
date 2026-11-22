`slog.NewLogLogger` takes the handler and the level its records get, and returns a `*log.Logger`: exactly what `accessLog` must return.

Give it `h` and `slog.LevelInfo`, return its result in place of the logger made with `log.New`, and remove the `"os"` import once nothing uses it.

If the test finds the lines on standard error, `accessLog` still returns the logger made with `log.New`. If the level is wrong, check the second argument: access lines are INFO.
