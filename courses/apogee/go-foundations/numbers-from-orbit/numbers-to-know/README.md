---
slug: numbers-to-know
title: Numbers to Know
type: code
position: 0
difficulty: 2
topics: printing, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Ingrid hands you the onboarding sheet, where Pipit-1's altitude appears three times, typed three ways: 525, 525km and 552. A number you name once can't disagree with itself.

`:=` declares a variable: it makes a new name and gives it a value. After that, using the name gives you its value:

```go
station := "Ystmark"
fmt.Println("Station:", station)
```

That prints `Station: Ystmark`. At Apogee, a number's name ends with its unit, as in `maskDeg` or `windKmh`.

## Your task

Below the given lines, declare `altitudeKm` with the value 525, then print it after the label `Altitude km:`, so the program prints:

```text
Pipit-1
Altitude km: 525
```
