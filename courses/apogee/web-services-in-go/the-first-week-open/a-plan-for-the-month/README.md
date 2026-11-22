---
slug: a-plan-for-the-month
title: A Plan for the Month
type: code_tests
position: 7
difficulty: 9
topics: errors, maps, time
created: 2026-11-13T06:35:19Z
modified: 2026-11-13T06:35:19Z
---
Friday, 11:00. Yusra: "TL;DR: Brinewatch has made 11 of the 12 bookings its plan allows in February, and asks what number 13 gets." A customer's plan is how many bookings it may make in a calendar month, counted in UTC and started again at midnight on the 1st.

`Quota` keeps the count in `usage.json`, so a restart forgets nothing, and calls `take` for every booking. Past the plan the answer is 403, not 429: Brinewatch's client waits as `Retry-After` says, a minute at most, so a 429 would have it knocking every minute until March.

```go
t := time.Date(2028, 12, 1, 0, 0, 0, 0, time.UTC)
fmt.Println(t.AddDate(0, 1, 0))
// 2029-01-01 00:00:00 +0000 UTC
```

## Your task

`take` counts nothing. Write it:

1. When `now`, in UTC, falls in a month other than `u.Month`, make `u.Month` that month's first instant and `u.Used` a new, empty map.
2. If `customer` has made `plan` bookings, return an `*APIError`: status 403, code `quota_exceeded`, message `plan of 12 bookings a month used up: the count starts again at 2028-03-01T00:00:00Z`, with its own plan and next month in RFC 3339.
3. Otherwise count the booking.
