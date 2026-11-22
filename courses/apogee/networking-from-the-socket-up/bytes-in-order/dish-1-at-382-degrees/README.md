---
slug: dish-1-at-382-degrees
title: Dish 1 at 382 Degrees
type: code
position: 0
difficulty: 7
topics: bytes, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Monday 27 March. Dish 1's new controller is in its radome, and its first status frame, read the way page 14 of the manual says, puts the dish at 382.15° in a 3,277 km/h wind. Sunniva: "Dish 1 at 382 degrees. Wind 3277. Neither true."

A number of two bytes or more can be stored either way round. Big-endian puts the most significant byte first, so `47 95` is 0x4795, 18325; little-endian puts it last, so the same bytes are 0x9547, 38215. `binary.BigEndian` and `binary.LittleEndian` read bytes in their own order:

```go
b := []byte{0x01, 0x2c}
fmt.Println(binary.BigEndian.Uint16(b))
// 300
fmt.Println(binary.LittleEndian.Uint16(b))
// 11265
```

## Your task

`main` reads Monday's status body little-endian, as page 14 says. Like every field of the protocol, these are big-endian: read them that way, so the program prints:

```text
az 183.25
el 41.70
wind 64.0
```
