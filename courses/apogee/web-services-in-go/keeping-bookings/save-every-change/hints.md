`f.mem.List` with `context.Background()` gives every booking in order, and `json.Marshal` turns that slice into one JSON array, with the names that `Booking`'s tags give its fields.

Hand `f.path` and the bytes to `replaceFile` and return what it returns. `Create` and `Delete` undo their change only when `save` returns an error, so every error has to come back.

If the read-only folder test fails, `save` returned nil although `replaceFile` failed, or it dropped an error from `List` or `json.Marshal` on the way.
