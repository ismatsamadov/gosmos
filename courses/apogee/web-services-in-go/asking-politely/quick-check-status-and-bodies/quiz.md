## The weather service answers `503 Service Unavailable`. What is `err` after `Do`?

- [x] `nil`: a 503 is still an answer
  > A 503 is an answer, so `Do` returns it with a nil error; the status is for the caller to check.
- [ ] An error that reads `503 Service Unavailable`
  > That error comes from the status check after `Do`, not from `Do` itself.
- [ ] An error only when the 503 has no body
  > A body, or no body, makes no difference: `Do` got an answer, so its error is nil.

## Nothing reads the 503's body. Must it still be closed?

- [ ] No: a body that is never read closes itself
  > A body never closes itself: closing it is always the caller's job.
- [x] Yes: until it is closed, its connection stays busy
  > Read or not, every body must be closed: `Close` is what lets its connection go.
- [ ] Only when the status is 200
  > Every body holds a connection, whatever the status above it says.

## When does the deferred `resp.Body.Close()` above run?

- [ ] Right away, on the line where it is written
  > `defer` only registers the call; it runs when the function returns.
- [ ] Only when the function reaches its last line
  > It runs on every return after the `defer` line, early ones included.
- [x] On every return after the `defer`, the 503's included
  > A deferred call runs whenever the function returns, so the early return closes the body too.
