---
slug: an-index-for-bookings
title: An Index for Bookings
type: code_tests
position: 4
difficulty: 9
topics: maps, profiling
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Thursday afternoon, the profile's next line: for each predicted pass, passplan read all 2,000 bookings of the next fourteen days, the passes customers have reserved, looking for the one that books it.

A search inside a loop does the whole search every time round. An index does it once: one walk through the bookings fills a map, and each pass then costs one lookup. The key holds what a booking and its pass share: the satellite, the station and the second the pass starts, its AOS, in Unix seconds, never a `time.Time`, whose `==` compares the time zone too:

```go
aos := time.Date(2029, 1, 16, 5, 29, 44, 0,
	time.UTC)
lan := aos.In(time.FixedZone("+0530", 19800))
fmt.Println(aos == lan, aos.Unix() == lan.Unix())
// false true
```

## Your task

Make `attach` build a `map[bookingKey]string` of booking IDs from `bookings` once, keeping the first when two bookings share a key, then give each pass the ID its own key finds.

> [!NOTE]
> A pass nobody booked finds `""` in the map, which is what `Planned` wants for it.
