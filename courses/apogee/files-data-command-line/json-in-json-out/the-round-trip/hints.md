Decode into a `Booking` first, and return nil with any decoding error.

A missing, empty or `null` ID all leave `ID` as the empty string, so one comparison catches the three. Return nil with an error made by `errors.New` for them.

Only a booking that passed the check is encoded: return what `json.Marshal` gives for it. The tags and `omitempty` in `booking.go` take care of the key names and the empty note.
