---
slug: refuse-misspelled-fields
title: Refuse Misspelled Fields
type: code_tests
position: 1
difficulty: 5
topics: errors, streams, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz wants tomorrow's passes to climb at least 10° above the horizon, a minimum elevation she typed as `min_elev_dg`. gctl read the settings file without a word and planned a pass that never climbs past 4°.

`Decode` skips any JSON field that matches no struct field. Call the decoder's `DisallowUnknownFields` before `Decode`, and such a field becomes an error that names it.

```go
in := `{"wind_kmh":64,"gust_kmh":81}`
dec := json.NewDecoder(strings.NewReader(in))
dec.DisallowUnknownFields()
var w struct{ WindKmh int `json:"wind_kmh"` }
err := dec.Decode(&w)
fmt.Println(err) // json: unknown field "gust_kmh"
```

## Your task

Implement `strictConfig`: decode one `Settings` value from `r` with a decoder that refuses unknown fields. Return the settings, or empty settings and the decoding error.
