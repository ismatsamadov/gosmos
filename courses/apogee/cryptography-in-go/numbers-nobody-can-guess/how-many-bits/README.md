---
slug: how-many-bits
title: How Many Bits
type: code
position: 4
difficulty: 6
topics: conversions, crypto
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Ingrid remembers when the ops console signed people in with six-digit codes, and she wants one line for her checklist: "Six digits. Basically as good as the new tokens? Yes or no."

A secret's strength is the number of guesses it takes, counted in bits: every bit doubles the guesses. A secret of `length` symbols, each picked at random from an alphabet of `alphabet` symbols, has length × log2(alphabet) bits, where log2 of a number is the power of 2 that makes it. `math.Log2` works it out:

```go
fmt.Println(math.Log2(32)) // 5
// three letters, A to Z:
fmt.Printf("%.1f\n", 3*math.Log2(26)) // 14.1
```

## Your task

`bits` multiplies the length by the alphabet's size. Make it return the length times `math.Log2` of the alphabet's size, so that `main` prints:

```text
6 digits: 19.9 bits
4 bytes in hex: 32.0 bits
rand.Text: 130.0 bits
32-byte key: 256.0 bits
```

> [!NOTE]
> Apogee's floor for any secret is 128 bits. Six digits are a million guesses; 128 bits are 2^128^.
