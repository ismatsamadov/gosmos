Call `checkShape` and then `checkPass` right after `decodeJSON`, the way `decodeJSON` is called: when one returns an error, hand it to `writeError` and return, since each check already says what the caller gets.

`errors.Is` tells an error that wraps `ErrConflict` apart from the store's other errors. Turn it into the 409 with the error's own text, which already reads `furrow-a over SQT at ...: already booked`, the way `getBooking` turns `ErrNotFound` into a 404.

If the 201 has no `Location`, set the header before `writeJSON`, which sends the status and every header with it: a header set afterward is never sent. If booking Monday's pass twice gets 500, the conflict hasn't become an `*APIError`.
