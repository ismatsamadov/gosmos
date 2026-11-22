`slices.Contains` reports whether a slice holds a value, and a customer's satellites are in its `SatIDs`. In each handler, `customerFrom` with the request's context gives the caller; `authenticate` has already made sure there is one.

In the two lists, keep a pass or a booking only when the caller owns its `SatID`, and an empty result must still be `[]`. In `getPass`, a pass that isn't the caller's doesn't match, so the loop ends in the same 404 as a missing one.

For `getBooking` and `cancelBooking`, get the booking from the store first, and when it isn't the caller's, answer what a missing one gets, `ErrNotFound` wrapped with `fmt.Errorf` as `booking bk-4510: not found`, before anything is deleted; one helper for both keeps them alike, and `bookings.go` needs `fmt`. In `createBooking`, the 403 is an `*APIError` with `http.StatusForbidden` and the code `forbidden`, after `checkShape` and before `checkPass`: if a missing Cirrine-1 pass gets 422, it comes too late.
