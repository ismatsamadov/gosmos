---
slug: bytes-are-not-letters
title: Bytes Are Not Letters
type: code
position: 3
difficulty: 5
topics: runes, strings
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Paz sends a voice memo: the Board thinks "Quebrada Añil", the road up to Salar Quieto, is too long for its column and cuts it a letter short. It measures names with `len`.

A Go string is a sequence of bytes, and a string literal stores its text as UTF-8. `len` counts bytes, and `s[i]` reads a single byte. ASCII characters take one byte each, but `ñ` and `°` take two. `utf8.RuneCountInString`, from the `unicode/utf8` package, counts runes instead: Unicode code points, one per character you see here:

```go
s := "45 °C"
fmt.Println(len(s))                    // 6
fmt.Println(utf8.RuneCountInString(s)) // 5
```

## Your task

Make the program print each name's length in runes: 7, 13 and 5.
