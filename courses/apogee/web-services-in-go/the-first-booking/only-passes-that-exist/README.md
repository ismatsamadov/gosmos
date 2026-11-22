---
slug: only-passes-that-exist
title: Only Passes That Exist
type: code_tests
position: 7
difficulty: 7
topics: errors, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
A Furrowlight test booked Thistle, the student satellite that fell out of orbit in July, and a pass that had begun an hour before. The service kept both.

A well-shaped booking must still match the service's own data. Compare times with `Equal`, which asks whether they are the same instant; `==` also compares how each was written:

```go
a, _ := time.Parse(time.RFC3339,
	"2027-12-06T15:02:10+01:00")
fmt.Println(a == a.UTC(), a.Equal(a.UTC()))
// false true
```

## Your task

`checkPass` returns nil. Return `invalidBooking` with the first that applies:

```go
"unknown satellite %q" // knownSat says no
"unknown station %q"   // knownStation says no
"no predicted pass of %s over %s at %s"
"bookings for this pass closed at %s"
```

The third is for a booking that matches no catalog pass in satellite, station, start and end; it names the first three. Bookings close ten minutes before the start, which ops calls AOS, because Ingrid's checklist must be done by then: from that very second, the pass is closed. Write times in UTC as RFC 3339.
