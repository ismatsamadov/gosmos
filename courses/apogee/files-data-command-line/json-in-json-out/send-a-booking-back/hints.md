`json.Marshal` already returns a byte slice and an error, the same two results `encodeBooking` returns.

Do not build the JSON by hand or drop the error from `Marshal`: a note with quotes must still decode to the same text, and the year 20270 must come back as an error with no bytes.
