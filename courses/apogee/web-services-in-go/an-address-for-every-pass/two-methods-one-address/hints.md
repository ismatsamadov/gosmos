Register the new route on the same path as `GET /bookings/{id}`, with `DELETE` in front of it. The two routes live side by side, and the mux picks one by the method.

`getBooking`, in `bookings.go`, shows the shape: read `id` with `PathValue`, call the store, and tell `ErrNotFound` apart from any other error with `errors.Is`. Any other error gets 500 `internal error`.

If the first cancellation gets 200, nothing wrote a status: write 204 with `w.WriteHeader` and write no body. If the second cancellation gets 204 too, the error from `Delete` is being dropped.
