Turn the `if` in `writeError` into a `switch` with no expression: the `*APIError` case first, then one case for each sentinel, each setting the status, the code and `err.Error()` as the message. Anything else keeps the 500.

`errors.Is(err, ErrNotFound)` is true for `booking bk-5499: not found`, because the store wrapped `ErrNotFound` with `%w`; comparing with `==` would miss it.

In `bookings.go`, delete the two `if errors.Is` blocks that build an `*APIError` from a store error, and the `errors` import with them, so each handler hands the store's error straight to `writeError`.
