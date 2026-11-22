---
slug: the-day-in-one-zone
title: The Day in One Zone
type: code_tests
position: 8
difficulty: 7
topics: slices, sorting, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz's plan for Thursday the 23rd and the network's list for the same day both hold four bookings, and only two of them match. One rises at 23:30 on the van's clock, when it is already Friday in UTC: the two Thursdays begin seven hours apart.

Which date a pass falls on depends on the clock reading it, so name the clock. Show each AOS in the chosen zone with `In`, write its date with `time.DateOnly`, the layout `2006-01-02`, and compare that with the day. The laptop's own zone never decides.

```go
van := time.FixedZone("", -7*60*60)
aos := time.Date(2027, time.September, 24,
	6, 30, 0, 0, time.UTC)
fmt.Println(aos.Weekday())         // Friday
fmt.Println(aos.In(van).Weekday()) // Thursday
```

## Your task

Implement `dayBookings`: return the bookings whose AOS falls on `day`, such as `2027-09-23`, on `zone`'s clock, ordered by AOS. Collect them in a new slice and leave the caller's slice as it was.
