---
slug: one-object-means-one
title: One Object Means One
type: code_tests
position: 2
difficulty: 5
topics: errors, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Nkem's review of the settings code reaches the van with one question: "What happens when the file holds two objects?" Today gctl uses the first and never mentions the second.

`Decode` reads one value and stops, whatever follows. To be sure the file held one object, ask for one more token: `dec.Token()` returns `io.EOF` only when nothing but white space is left.

```go
in := `{"wind_kmh":64} {"wind_kmh":81}`
dec := json.NewDecoder(strings.NewReader(in))
var first map[string]int
_ = dec.Decode(&first)
_, err := dec.Token()
fmt.Println(err == io.EOF) // false
```

## Your task

`decodeOne` already decodes the settings strictly. After the object, accept only `io.EOF`: for anything else, a second value or stray text, return empty settings and `ErrTrailing`. Trailing white space is fine.
