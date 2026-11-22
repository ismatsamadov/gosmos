---
slug: one-chunks-size
title: One Chunk's Size
type: code_tests
position: 1
difficulty: 9
topics: conversions, http, strings
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Monday afternoon. Nkem asks how the recorder sends a pass whose length it learns only when the satellite sets. Go's client streams it in chunks, each starting with a line giving its size in hexadecimal: `c` for 12 bytes, `1000` for 4,096.

The size is hex digits of either case, and the station server takes 1 to 16 of them. Spaces or tabs, then a `;` and an extension the server ignores, may follow, as in `400;name=value`. `strconv.ParseInt` reads base 16, but also takes a sign, so the digits need checking first:

```go
n, err := strconv.ParseInt("+1f", 16, 64)
fmt.Println(n, err) // 31 <nil>
```

## Your task

`chunkSize` reads the line as decimal with `strconv.Atoi`. Cut it at the first `;`, trim spaces and tabs from the end of what is left, require 1 to 16 hex digits and parse them in base 16. Refuse anything else with 400 and the line quoted:

```text
bad chunk size "zz"
```
