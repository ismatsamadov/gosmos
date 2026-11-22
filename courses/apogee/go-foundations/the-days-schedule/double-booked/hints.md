When the outer loop is at pass `i`, which passes has it already been paired with, and which are left?

Start the inner loop at `i + 1`, so no pass meets itself and no pair comes up twice. Check that the stations match, then compare the times with `Before`.

The time check has two halves joined by `&&`: `a.AOS.Before(___)` and `b.AOS.Before(___)`. For touching passes, one half is false.
