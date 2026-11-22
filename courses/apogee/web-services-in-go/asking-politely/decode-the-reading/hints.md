Decode into a `Weather` variable with `json.NewDecoder(r)` and its `Decode` method, and return the decoder's error as it is: broken JSON and an empty body both come back from there.

A field missing from the JSON stays at its zero value: an empty `Station`, and an `At` whose `IsZero` method reports true. Check both after decoding and return an error of your own made with `errors.New`.

If the reading with `humidity_pct` fails with `unknown field`, you asked the decoder to refuse unknown fields: leave that out, since this service adds fields whenever it likes.
