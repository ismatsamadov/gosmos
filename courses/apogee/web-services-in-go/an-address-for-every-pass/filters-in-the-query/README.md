---
slug: filters-in-the-query
title: Filters in the Query
type: code_tests
position: 4
difficulty: 6
topics: http, json, slices
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T16:20:02Z
---
Furrowlight's program asks for `/passes?sat=furrow-a&station=SQT` and gets back every pass the service offers its customers, its whole catalog, Pipit-1's and Brine-3's included.

`r.URL.Query().Get` returns `""` for a parameter that isn't there, so an empty value means no filter. A value the service doesn't know is the caller's mistake: answer 400 before listing anything. And start the list as an empty slice, because `encoding/json` writes a nil slice as `null`:

```go
var none []string
a, _ := json.Marshal(none)
b, _ := json.Marshal([]string{})
fmt.Println(string(a), string(b)) // null []
```

## Your task

`listPasses` answers with the whole catalog. Keep the passes that match `sat` and `station` when each is given, and answer `[]` when none match. Answer 400 with `http.Error` for a value that `knownSat` or `knownStation` doesn't know:

```text
unknown satellite "thistle"
unknown station "ZZZ"
```

> [!NOTE]
> Brine-3 never rises over Ystmark, so the answer to `?sat=brine-3&station=YST` is `[]`.
