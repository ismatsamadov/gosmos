---
slug: store-errors-right-statuses
title: Store Errors, Right Statuses
type: code_tests
position: 7
difficulty: 8
topics: errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Cirrine's program canceled a booking that Monday's restart had already lost, and got a 500. It woke the engineer who answers Cirrine's alarms that week, their on-call, for a booking that wasn't there.

`getBooking` and `createBooking` each turn a store error into a status themselves; `cancelBooking` did too, until Tuesday's tidy-up shortened it. `errors.Is` finds a sentinel error however deeply it is wrapped, so the mapping can live in `writeError`, which every error passes through:

```go
switch {
case errors.Is(err, fs.ErrNotExist):
	fmt.Println("no such report")
case errors.Is(err, fs.ErrPermission):
	fmt.Println("not yours to read")
}
```

## Your task

1. In `writeError`, after the `*APIError` case, answer `ErrNotFound` with 404 `not_found` and `ErrConflict` with 409 `conflict`, keeping the error's own message, such as `booking bk-5499: not found`.
2. Take the store errors' special cases out of `getBooking` and `createBooking`.

> [!NOTE]
> Keep the `*APIError` case first: an error that is both answers as its `*APIError` says.
