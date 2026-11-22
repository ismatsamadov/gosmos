Convert each count with `float64(...)` before you divide: between two `int` values, `/` keeps only the whole part, so 1,177 / 1,199 would be 0.

Check the counts before the division, booked first: for each one that can't be, return 0 and an error made with `errors.New` from the exact text in the task.

If the test reports NaN, the division ran with nothing booked, so the check has to come first. If an error's text differs, compare it with the task letter by letter: lowercase, no final period.
