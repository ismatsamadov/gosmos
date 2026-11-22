Make two bookings first, such as Furrow-A's and Furrow-B's passes over Salar Quieto on the same day, Furrow-A's the earlier. Each subtest calls `open()` for a store of its own, so every check starts empty.

`errors.Is` checks for `ErrConflict` or `ErrNotFound` however the store wrapped them. Report with `t.Errorf` inside the subtest, so a failure names its case; a `t.Fatalf` for a `Create` that should have worked stops only that subtest.

If a mutant survives, read it in `main_test.go`. The store that lists bookings in the order they were made sorts nothing, so it passes unless your subtest books the later pass before the earlier one.
