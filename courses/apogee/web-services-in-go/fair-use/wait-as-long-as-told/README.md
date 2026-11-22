---
slug: wait-as-long-as-told
title: Wait as Long as Told
type: code_tests
position: 4
difficulty: 8
topics: conversions, http, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra is sending Brinewatch a corrected sample client, and its retry loop must wait exactly as long as the service says. Brinewatch's engineers will copy it line for line.

A client reads `Retry-After` from a 429 and waits that many seconds before trying again. The booking service always sends a whole number, but a careful client never trusts a header blindly: a missing or unreadable one means 1 second, and nothing waits more than 60:

```go
n, err := strconv.Atoi("soon")
fmt.Println(n, err != nil) // 0 true
n, err = strconv.Atoi("-5")
fmt.Println(n, err)        // -5 <nil>
```

## Your task

`retryDelay` returns 0. Read `Retry-After` from `res`'s header as whole seconds with `strconv.Atoi`, and return that many seconds. Return 1 second when the header is missing, unreadable or below 1, and 60 seconds for anything above 60.

> [!WARNING]
> Compare the number with 60 before you make it a `time.Duration`: enough seconds overflow one.
