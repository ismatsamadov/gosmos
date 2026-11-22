---
slug: a-week-at-a-time
title: A Week at a Time
type: code_tests
position: 5
difficulty: 7
topics: http, time
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Yusra's TL;DR: "Furrowlight wants to book a week at a time." Hollis's program will pick when each week starts.

Query values are text: a time needs `time.Parse`, a count `strconv.Atoi`, and either can fail. A missing value takes its default; one that is unreadable or out of range is the caller's mistake:

```go
n, err := strconv.Atoi("seven")
fmt.Println(n, err != nil) // 0 true
t, err := time.Parse(time.RFC3339,
	"2027-11-29T09:00:00+01:00")
fmt.Println(t.UTC(), err)
// 2027-11-29 08:00:00 +0000 UTC <nil>
```

## Your task

`window` ignores the query and returns `now` and 14. Make it return `from`, an RFC 3339 time defaulting to `now`, and `days`, a whole number from 1 to 14 defaulting to 7, or one of these errors:

```text
from must be an RFC 3339 time
days must be a whole number from 1 to 14
```

> [!NOTE]
> `listPasses` lists the passes starting at or after `from` and before `from` plus `days` days.
