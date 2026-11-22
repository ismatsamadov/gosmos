`token` with `furrowlightKey`, then with `brinewatchKey`, gives the two tokens; the keys are in `fixtures.go`. Decode `GET /passes` into a `[]PassView`, and build a `Booking` from the first pass's satellite, station, AOS and LOS.

Send that booking twice with `"Idempotency-Key", "rehearsal-1"` after the body, then compare the two statuses and the bodies with `bytes.Equal`. Decode the first body into a `Booking` to learn its ID.

Check both the status and the `errorCode()` of each refusal, so an error in plain text fails too. For the limit, send `GET /bookings` as Brinewatch until a 429 comes or 11 requests have gone: a service that never limits never sends one.
