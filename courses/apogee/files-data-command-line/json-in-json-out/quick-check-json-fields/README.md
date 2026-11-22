---
slug: quick-check-json-fields
title: Quick Check: JSON Fields
type: quiz
position: 3
difficulty: 4
topics: structs
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Before the van's replies go out, Nkem reviews the exchange format at the next signal, starting with this version of the type:

```go
type Booking struct {
    ID      string `json:"id"`
    SatID   string `json:"sat_id"`
    station string
}
```

Here `station` holds `"SQT"`. Encoding and decoding see only what the json package can reach: exported fields, under their tags.

## Your task

Answer Nkem's three questions: which key `SatID` gets, what happens to `station`, and what `json.Unmarshal` returns for input that stops halfway.
