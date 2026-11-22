Loop over the bookings and keep the ones whose AOS, shown in `zone` with `In` and written with the `time.DateOnly` layout, equals `day`. Append them to a new slice.

Sort only the new slice, with `slices.SortFunc` and the times' own `Compare` method. Sorting the argument would reorder the caller's bookings.

Formatting the AOS without `In` writes its UTC date, and `Local` writes the laptop's: both put a pass at 23:30 on the van's clock on the wrong day.
