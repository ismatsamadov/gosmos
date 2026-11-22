---
slug: every-answer-has-a-status
title: Every Answer Has a Status
type: code_tests
position: 5
difficulty: 5
topics: http, switch
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra is writing the pages that tell customers' engineers how to use the service, the customer guide. Her TL;DR: "One line for every kind of answer a program can get, so nobody has to guess."

Every answer carries a three-digit status code, and its first digit is its class: 2xx means the request worked, 4xx that the request was wrong, 5xx that the server failed. One 4xx is different: `429 Too Many Requests` means the request was fine but came too soon, so wait and send it again.

```go
fmt.Println(http.StatusNotFound) // 404
fmt.Println(http.StatusText(503))
// Service Unavailable
```

## Your task

`advice` returns `""`. Make it return the line for each status code:

```text
200 to 299        done
429               wait and retry
other 400 to 499  fix the request
500 to 599        retry later
anything else     unexpected
```

> [!WARNING]
> 200 is `done` too: a check such as `code > 200` turns away the most common answer of all.
