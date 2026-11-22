---
slug: taking-a-line-apart
title: Taking a Line Apart
type: code
position: 0
difficulty: 4
topics: strings
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Thistle's telemetry comes down as text lines: satellite, sequence number and time, then `key=value` fields. Tonight a line arrived with a doubled space, and the parser printed an empty sequence number.

`strings.Split(s, " ")` cuts at every single space, so two spaces in a row leave an empty field between them. `strings.Fields(s)` splits around runs of white space (spaces, tabs, newlines) instead, and never returns an empty field:

```go
split := strings.Split("YST  PEL", " ")
fields := strings.Fields(" YST  PEL ")
fmt.Println(len(split))  // 3: "YST", "", "PEL"
fmt.Println(len(fields)) // 2: "YST", "PEL"
```

## Your task

In `main`, split `line` with `strings.Fields` instead of `strings.Split`, so the program prints six fields, the satellite and the sequence number.
