---
slug: twelve-bytes-every-time
title: Twelve Bytes Every Time
type: code_tests
position: 6
difficulty: 5
topics: arrays, bytes, types
created: 2026-11-04T17:00:51Z
modified: 2026-11-05T08:05:37Z
---
Every Pipit broadcasts a short status message, its beacon, in the same layout of exactly twelve bytes. Paz wants the decoder to express that size in its type, and needs a recorded Pipit-1 beacon copied with Pipit-2's number for tests, leaving the recording as it was.

`[12]byte` is an array: twelve bytes, with the length part of its type. `[11]byte` is a different type. Arrays are values: assignment and function parameters copy every element. Arrays of comparable elements can be compared with `==`:

```go
code := [3]byte{'P', 'E', 'L'}
draft := code
draft[0] = 'Y'
fmt.Println(code == draft) // false
```

## Your task

Complete `sameBeacon` using array equality. In `withSatellite`, set byte 2 of the parameter, the satellite number, to `number` and return the array. The caller's original must stay unchanged, including its byte 2.
