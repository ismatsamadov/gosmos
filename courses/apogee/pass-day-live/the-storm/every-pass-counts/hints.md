Build the report from what the run already holds: the bookings it answers for, the unplanned ones with their reasons, and the outcome of each matched booking's contact. A booking at a station the rehearsal doesn't fly gets no row.

Per station, booked is every booking there and lost every one not flown, whether unplanned, failed or skipped. `time.Time.ISOWeek` gives the title, and the budget stays in integers, `booked * (10000 - target) / 10000`, so 210 bookings at 9950 allow 1.

Print `-` with `%6s` when nothing is booked, sort the lost bookings by AOS then ID, and write each one's outcome as its LOS line says it, or `unplanned` and the reason, before the verdict. Compare your `week.txt` line by line with `week.golden` in `testdata`: the first line that differs is in the test's message.
