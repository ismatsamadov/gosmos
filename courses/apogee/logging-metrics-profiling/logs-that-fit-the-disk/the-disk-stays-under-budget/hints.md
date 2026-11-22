Make the rotator as a struct literal with its `Dir`, `Program`, `MaxBytes` and `Now`, then set its `OnOpen` to a function that calls `prune` for the same folder and program, with `now()`, `keepDays` and `totalBytes`. Pruning that fails must not stop the log, so the function drops `prune`'s error. Then call `r.open()`, and return its error with two nils if it fails.

Build the handler from the inside out: `slog.NewJSONHandler` on the rotator with no options, given to `NewSampler`, given to `stamp` with `now`. Return `slog.New` of that, the rotator and nil, and remove the `path/filepath` import once nothing uses it.

If the test can't find the warnings that open each minute, the records carry the Mac's clock or the sampler counted minutes by it: `stamp` must be the outermost handler. If 19 November's file is still there, nothing prunes the folder: give the rotator its `OnOpen` before you open it.
