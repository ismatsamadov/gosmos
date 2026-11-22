---
slug: two-methods-one-address
title: Two Methods, One Address
type: code_tests
position: 6
difficulty: 7
topics: errors, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Brinewatch, the ship-tracking customer, changes its plans by email almost every week, and this week's email cancels `bk-4496`, Brine-3's pass on Friday. Yusra wants a cancellation to be one request instead of an email.

Two patterns can share a path when their methods differ, each with its own handler, and the mux lists every method in `Allow`. A request that removes something has nothing to send back, so it answers `204 No Content`: a status, and no body at all.

```go
mux.HandleFunc("GET /notes/{n}", readNote)
mux.HandleFunc("PUT /notes/{n}", saveNote)
// POST /notes/7: 405, Allow: GET, HEAD, PUT
```

## Your task

1. In `routes`, register `DELETE /bookings/{id}` for `s.cancelBooking`.
2. In `cancelBooking`, delete the booking with `s.Store.Delete(r.Context(), id)` and answer 204. For a booking it doesn't have, `Delete` returns an error wrapping `ErrNotFound`: answer it as `getBooking` does, with 404 and the error's text:

```text
booking bk-4496: not found
```
