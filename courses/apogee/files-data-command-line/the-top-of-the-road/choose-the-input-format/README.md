---
slug: choose-the-input-format
title: Choose the Input Format
type: code_tests
position: 1
difficulty: 7
topics: errors, streams, switch, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Yusra, who looks after Apogee's customers, starts with a one-line summary, as always: "TL;DR: Thursday's bookings come in two files, my export as JSON and Hollis's sheet as CSV. Same passes; read whichever one Paz names."

Choose the decoder from what the caller says, never by guessing from the bytes: an empty or cut-off file looks like neither format, and a guess fails differently on every file. A `switch` on the format's name sends each one to its reader and refuses the rest:

```go
switch unit {
case "km":
	return value, nil
case "m":
	return value / 1000, nil
default:
	return 0, fmt.Errorf("unknown unit %q", unit)
}
```

## Your task

Implement `readBookings`: for `"json"` return what `decodeBookings` returns, for `"csv"` what `readBookingSheet` returns, and for any other format an error that names it. Pass their errors back unchanged.
