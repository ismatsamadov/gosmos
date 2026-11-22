---
slug: your-data-has-landed
title: Your Data Has Landed
type: code_tests
position: 9
difficulty: 9
topics: context, http, structs
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday 21 January, 14:42 UTC: Furrow-A's image from today's pass over Salar Quieto has landed, and Furrowlight's receiver should hear about it within a minute.

The notifier turns a landing into an event for the right customer: it looks up the booking, works out whose satellite flew the pass, and hands the sender one `Delivery`. Event IDs count up from `evt-000001`, and an event keeps its ID on every attempt, so a receiver can drop a copy it has already seen:

```go
var n atomic.Int64
fmt.Printf("evt-%06d\n", n.Add(1))
// evt-000001
```

## Your task

`Landed` does nothing. Get the booking from `n.Store` and its customer from `n.owner`, build a `data.landed` event with the next ID from `n.events`, the booking's ID, satellite and station, and the landing's file, bytes and time, then deliver it with `n.Sender`. Return any error.
