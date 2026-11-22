Decode into a `[]Booking` variable with `json.Unmarshal`, passing its address, and return nil with any error, even when some entries were already decoded.

After a successful decode, a nil slice means the export said `null`. Replace it with `[]Booking{}` then: making the slice before decoding is not enough, because `null` sets it back to nil.
