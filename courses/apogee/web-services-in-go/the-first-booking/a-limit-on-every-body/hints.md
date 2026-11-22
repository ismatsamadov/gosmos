Wrap the body where the decoder is made, giving `http.MaxBytesReader` the writer, `r.Body` and `maxBodyBytes`. The decoder then reads through the limit without knowing it is there.

Both `Decode` calls send their errors to `bodyError`, so one check there covers both. `errors.AsType`, given the limit's error type in brackets, reports whether `err` is that error; check it before the empty-body case.

If one byte too many gets 400 `request body: http: request body too large`, the limit works but `bodyError` never checks the error's type. If 70,000 spaces after the booking get 400, the second `Decode`'s error isn't going through that check.
