Make a decoder for `r.Body` with `json.NewDecoder` and call its `Decode` method with `&b`, the address of the booking, so it can fill the fields in.

When `Decode` returns an error, answer it with `http.Error` and `http.StatusBadRequest`, and return before `Create` runs. `Create` returns the booking with its ID and an error; the 201 carries that booking.

If the broken bodies get 201, or the store holds more than Yusra's bookings afterward, a failed decode still reaches `Create`. If the 201 has no `id`, you are answering the booking from before `Create`.
