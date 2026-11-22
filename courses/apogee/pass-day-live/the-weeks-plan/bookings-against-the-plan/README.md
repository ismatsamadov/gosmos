---
slug: bookings-against-the-plan
title: Bookings Against the Plan
type: code_tests
position: 2
difficulty: 9
topics: json, maps, sorting
created: 2026-11-10T20:36:24Z
modified: 2026-11-12T02:35:45Z
---
Tuesday 3 July. Yusra has exported the week's 1,495 bookings from bookingd, the booking service, and a booked pass with no dish is lost before the week begins.

`bookings.json` is a JSON array of `Booking` (`types.go`), read strictly; no file means no bookings. A booking names no pass: it matches the contact with its satellite, its station and its AOS, to the second. Bookings outside the window, or at a station the rehearsal doesn't fly, are left out.

## Your task

Save `reports/unplanned.txt`, even empty: a line per unmatched booking, by AOS then ID, with the first of these reasons that fits:

```text
unknown satellite
unknown station
no such pass
not servable
too low
no dish free
```

In Yusra's export, two bookings don't match, Hollis's among them:

```text
bk-9902 Thu 04:58:30 furrow-a/CWR no such pass
```

Then print the count. Once she removes both, the week's export gives:

```text
bookings: 1493, unplanned 0
```

A broken export is an error that names the file, such as `bookings.json: json: unknown field "sat"`.

> [!TIP]
> An unknown station is in no network at all, so it can't be one the rehearsal doesn't fly: it stays, as unplanned.
