---
slug: the-right-person
title: The Right Person
type: code_tests
position: 5
difficulty: 8
topics: conditionals, metrics, time
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Six of the eleven pages were single failed passes that could have waited for daylight, and each one woke Tomás. Ingrid: "Wake the person on call for what can't wait. The rest goes on my list."

A rule's severity says who reads its alerts. Critical pages the person on call at the alert's time: the latest shift on the rota that started at or before it, so at 07:59 on a Monday last week's name still holds. Warning goes on the morning list, Ingrid's list for the 08:00 meeting, and Info is only logged:

```go
mon := time.Date(2029, 1, 1, 8, 0, 0, 0, time.UTC)
at := mon.Add(-time.Minute) // Monday 07:59
fmt.Println(!mon.After(at)) // false: not yet
```

## Your task

`route` leaves `a.To` empty. Set it by `sev`: for `Critical`, the `Name` of the latest shift in `rota` whose `From` is at or before `a.At`; for `Warning`, `"morning list"`; for `Info`, `"log"`. Return the alert, and leave `rota` as it was.

> [!NOTE]
> A RESOLVED is routed the same way: it goes to whoever is on call at its own time, and wakes nobody.
