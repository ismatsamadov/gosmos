The failing test shows the right text but compares values. Which value should `accept` hand back instead of a new error?

In `shouldPage`, combine two conditions with `&&`: `err` isn't `nil`, and `errors.Is` says it isn't `ErrDuplicate`.
