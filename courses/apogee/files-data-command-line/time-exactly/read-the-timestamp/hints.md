The layout is the first argument to `time.Parse` and the text the second. `time.RFC3339` is a constant in the `time` package; you don't write the layout by hand.

`time.Parse` already returns a time and an error, in the order `parseStart` returns them, and the zero time when it fails.

A layout written by hand that ends in a literal `Z` refuses `-07:00`. The `RFC3339` layout ends in `Z07:00`, which accepts `Z` or an offset.
