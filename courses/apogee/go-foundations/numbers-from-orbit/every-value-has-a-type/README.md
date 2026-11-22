---
slug: every-value-has-a-type
title: Every Value Has a Type
type: code
position: 1
difficulty: 2
topics: types, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
The sheet's four facts about Pipit-1 are a name, a whole number, a decimal and a yes or no. Its radio frequency is in MHz: millions of cycles per second. In Go, every value has a type, and a variable keeps its type for good.

`var` declares a variable with a name, a type and a value. Four of Go's basic types are `string` (text), `int` (whole numbers), `float64` (numbers that can have a fraction) and `bool` (`true` or `false`):

```go
var station string = "Lantern Atoll"
var altM int = 3
var dishM float64 = 3.7
var hasXBand bool = false
```

## Your task

Below `name`, declare `altitudeKm` (525), `beaconMHz` (437.25) and `deployed` (true) with `var` and the right types. Then print each after its label, so the program prints:

```text
Name: Pipit-1
Altitude km: 525
Beacon MHz: 437.25
Deployed: true
```
