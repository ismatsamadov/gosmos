---
slug: a-404-is-an-answer
title: A 404 Is an Answer
type: code_tests
position: 3
difficulty: 5
topics: errors, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
This morning a typo asked the weather service for `CRW` instead of Cape Wren's `CWR`, and the poll reported `invalid character 'o'`. Its decoder had read the plain-text answer `no station "CRW"` as if it were JSON.

`Do` returns an error when no answer comes back at all. A 404 or a 503 is an answer: `err` is nil, and the status sits in the response, as a number and as text.

```go
fmt.Println(resp.StatusCode) // 404
fmt.Println(resp.Status)     // 404 Not Found
```

## Your task

`Reading` decodes whatever comes back. Before decoding, return an error unless the status is `http.StatusOK`, naming the station and the status:

```text
weather for CRW: 404 Not Found
```

> [!NOTE]
> That early return leaves the answer's body open. Closing every body is the next lesson.
