## If `parsePassID` returns an error, does the second `if` still run?

- [ ] No: `t.Errorf` ends the test at once
  > Ending the test is what `t.Fatalf` does; `t.Errorf` records the failure and returns.
- [x] Yes: `t.Errorf` marks the test failed, and the test carries on
  > The test keeps going, so the next check runs on an `id` that was never parsed.
- [ ] Only when the test runs with `go test -v`
  > `-v` changes how much `go test` prints, not which lines of a test run.

## Which call fits the first `if` better?

- [x] `t.Fatalf`, because nothing after it makes sense without a parsed ID
  > Stopping there spares you a second, misleading failure about a station that was never read.
- [ ] `t.Errorf`, because a test should report every problem it finds
  > Reporting every problem suits checks that stand on their own, but the station check depends on the parse.
- [ ] `panic`, so the failure can't be missed
  > A panic ends the whole test run, not only this test, and prints a stack trace instead of a clear failure.
