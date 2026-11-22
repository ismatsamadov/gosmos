Call `DisallowUnknownFields` on the decoder before the first `Decode`. From then on, a field that `Booking` doesn't have makes `Decode` return an error that names it.

After the first value, call `Decode` again into an empty struct. Only `io.EOF` means nothing but white space followed; anything else, an error or a second value, gets the message about unexpected data.

If the empty body says `request body: EOF`, `bodyError` doesn't yet turn `io.EOF` into the empty-body message. If two bookings in one body pass, the second `Decode` is missing or its result is ignored.
