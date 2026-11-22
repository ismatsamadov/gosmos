## The body says `"sat"` where `Booking` has `sat_id`. What does `err` hold?

- [ ] Nothing: `err` is nil and `b.SatID` stays empty
  > That is what a decoder does without `DisallowUnknownFields`, the very bug Furrowlight's booking hit.
- [x] An error that names the unknown field, `"sat"`
  > `DisallowUnknownFields` turns any field that `Booking` doesn't have into an error.
- [ ] Nothing: the decoder puts `furrow-a` in `b.SatID`, the closest field
  > The decoder matches names exactly, apart from letter case; it never guesses that `sat` meant `sat_id`.

## Why does the code call `Decode` a second time?

- [ ] To read the rest of a body too large for one call
  > One call reads a whole value, however large it is.
- [ ] To check that the first value was an object
  > The first `Decode` already fails when the value doesn't fit `b`.
- [x] One call reads one value, so only another call sees what follows it
  > `extra` is `io.EOF` when nothing but white space follows the booking.

## What does `errors.AsType[*APIError](wrapped)` return?

- [x] `apiErr` and true
  > `AsType` looks through every error wrapped with `%w` until it finds the type it asks for.
- [ ] A nil error and false, since `fmt.Errorf` made `wrapped`
  > With `%w`, `fmt.Errorf` keeps `apiErr` inside `wrapped`, where `AsType` finds it.
- [ ] A new `*APIError` whose message starts with `booking: `
  > It returns `apiErr` itself, message unchanged; only `wrapped`'s own text has the prefix.
