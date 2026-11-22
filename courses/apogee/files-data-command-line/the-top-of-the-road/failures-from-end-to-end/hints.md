Replace each `_` that throws an error away with `err`, and follow every call with a check that returns `Plan{}` and the wrapped error. `fmt.Errorf` with the `%w` verb puts words in front and keeps the original error inside.

`loadStations` reads the table built into gctl; its error can be returned as it is. Then hand the bookings and the stations to `validateBookings`, and return its error, wrapped, before `newPlan` runs.

Finish by returning what `newPlan` returns: both of its values, its own error included.
