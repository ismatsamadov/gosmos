---
slug: decode-the-reading
title: Decode the Reading
type: code_tests
position: 6
difficulty: 6
topics: json, streams
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Last spring the weather service started adding `humidity_pct` to its readings, and Tomás has more fields in mind. The poll must take what it needs and let the rest go by, without breaking.

`json.NewDecoder(r).Decode(&v)` reads one JSON value straight from a stream, such as a response body. It fills the fields `v` knows and skips the ones it doesn't, which is right for a service you don't own; what you need, check yourself afterwards.

```go
var s struct {
	Code string `json:"code"`
}
err := json.NewDecoder(strings.NewReader(
	`{"code":"PEL","gulls":3}`)).Decode(&s)
fmt.Println(s.Code, err) // PEL <nil>
```

## Your task

Write `decodeWeather`, which `Reading` calls with the body: decode one `Weather` from `r`, and return an error when the JSON is broken or when `station` or `at` is missing. Extra fields are fine.

> [!WARNING]
> Don't call `DisallowUnknownFields` here. The booking service will be strict about its own input; a client of someone else's service has to be lenient.
