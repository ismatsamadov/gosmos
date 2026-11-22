---
slug: point-dish-1
title: Point Dish 1
type: code_tests
position: 4
difficulty: 7
topics: bytes, conversions, errors
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tomás wants Dish 1 at azimuth 183.25°, elevation 41.70°, and then a STATUS, to see whether the dish ends up where the POINT command said. "Two numbers, four bytes. Let's see if the dish agrees!"

A POINT's body is the azimuth as a `uint16`, then the elevation as an `int16`, each in hundredths of a degree. Multiply by 100 and round with `math.Round` before converting: a conversion drops the fraction, and 0.29 × 100 is 28.999999999999996 in a `float64`. `binary.BigEndian.AppendUint16` appends the two bytes, most significant first:

```go
deg := 0.29
x := deg * 100
fmt.Println(int(x), int(math.Round(x)))
// 28 29
b := binary.BigEndian.AppendUint16(nil, 640)
fmt.Printf("% x\n", b) // 02 80
```

## Your task

`pointBody` returns `nil`. Round both angles to hundredths, refuse an azimuth outside 0 to 359.99 or an elevation outside 0 to 90.00 after rounding, and append the azimuth, then the elevation.
