Take the year and the month from `now.UTC()`. `time.Date` with day 1, a zero time of day and `time.UTC` gives that month's first instant, and `Equal` tells whether it is `u.Month`.

`AddDate(0, 1, 0)` on the first of a month gives the first of the next one, December's included, and `Format(time.RFC3339)` writes it as the message wants.

If February 2029 is refused, your month forgot its year; if a clock reading 21:30 at UTC−3 on 29 February is refused, it used the clock's own zone instead of UTC.
