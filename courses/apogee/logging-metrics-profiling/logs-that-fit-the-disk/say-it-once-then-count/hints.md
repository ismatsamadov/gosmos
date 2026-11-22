Begin with the level: a record at `slog.LevelError` or above goes to `s.inner.Handle` at once, uncounted. Every other record is counted under the lock, so lock `s.shared.mu`, ask `s.shared.tallyOf(r)` for its tally, and unlock as soon as you have decided.

Add one to the tally's `seen`, then keep the record when `seen` is at most 10 or leaves 10 when divided by 100. Before you change anything, note the tally's `dropped`: a kept record takes that number and sets `dropped` back to 0, while a dropped one adds one to `dropped` and to `s.shared.total` and returns nil.

If line 11 is the 100th warning rather than the 110th, check the remainder you compare with. If the first ten lines carry `"dropped":0`, add the attribute only when the number you noted is above 0.
