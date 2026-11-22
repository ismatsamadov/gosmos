---
slug: a-session-behind-the-cookie
title: A Session Behind the Cookie
type: code_tests
position: 1
difficulty: 9
topics: maps, mutexes, security, time
created: 2026-11-13T03:49:38Z
modified: 2026-11-13T12:18:54Z
---
In the draft, a session's ID is the customer's ID. Nkem reads it and asks one question: "What does the page show when Hollis edits his cookie to say `cirrine`?"

Cirrine's bookings: a browser's user can write anything into a cookie. So the cookie carries only a session ID, a random name nobody could guess, and the service keeps whose session it is and when it ends. `Text` from `crypto/rand` returns 26 random characters, `A` to `Z` and `2` to `7`, about 130 bits:

```go
fmt.Println(rand.Text())
// RX4KWQ2ZB7NDJ5MC3TVLHY6PAE
```

A session ends 30 minutes after it was last used: each successful lookup moves its end to `sessionLife` after `s.now()`.

## Your task

1. `Start` returns the customer's ID. Keep a `session` for `customer` in `s.byID` under a new `rand.Text()` ID, ending `sessionLife` after `s.now()`, and return the ID.
2. `Lookup` trusts any ID. While `s.now()` is before the session's end, move the end and return its customer and true. Delete a session that has ended, and return `"", false` for it and for an unknown ID.

Hold `s.mu` in both.
