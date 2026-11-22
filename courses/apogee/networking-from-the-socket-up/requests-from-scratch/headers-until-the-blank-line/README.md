---
slug: headers-until-the-blank-line
title: Headers Until the Blank Line
type: code_tests
position: 4
difficulty: 8
topics: http, maps, strings
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Headquarters' collector, the program at Port Ellery that will fetch each recording from the hut, says who it is and what it wants in its header fields. The station server must read every one.

Each header line is a name, a colon and a value, up to the blank line that ends the head. Split a line at its first colon with `strings.Cut`, keep the name in lower case, since names match in any letter case, and trim spaces and tabs around the value. A `Header`, a `map[string][]string`, keeps every value of a repeated name in order:

```go
key, value, ok := strings.Cut("batt_v=7.91", "=")
fmt.Println(key, value, ok) // batt_v 7.91 true
```

## Your task

`readHeader` in `request.go` returns an empty `Header`. Read lines with `readLine` until the empty one, split each at its first colon, lower-case the name, trim spaces and tabs around the value with `strings.Trim`, and append the value to the name's.
