---
slug: refuse-what-you-dont-know
title: Refuse What You Don't Know
type: code_tests
position: 2
difficulty: 6
topics: errors, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Furrowlight's program sent `"sat":"furrow-a"` instead of `"sat_id"`. The decoder skipped the field it didn't know, and the service tried to book a pass for no satellite at all.

A decoder ignores unknown fields unless you call its `DisallowUnknownFields` method first; then a misspelled field is an error. And `Decode` reads one JSON value and stops, so decode once more to see what follows: white space alone gives `io.EOF`.

```go
dec := json.NewDecoder(strings.NewReader(
	`{"code":"YST"} {"code":"PEL"}`))
var st Station
err := dec.Decode(&st)
fmt.Println(st.Code, err) // YST <nil>
err = dec.Decode(&st)
fmt.Println(st.Code, err)    // PEL <nil>
fmt.Println(dec.Decode(&st)) // EOF
```

## Your task

Make `decodeJSON` strict: refuse unknown fields, then decode again into `&struct{}{}` and refuse anything but `io.EOF`. In `bodyError`, say the body was empty when `err` is `io.EOF`. Each message starts with `request body: `, then:

```text
json: unknown field "sat"
unexpected data after the JSON value
empty
```
