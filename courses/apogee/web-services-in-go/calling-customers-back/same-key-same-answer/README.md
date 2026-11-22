---
slug: same-key-same-answer
title: Same Key, Same Answer
type: code_tests
position: 8
difficulty: 9
topics: http, maps, mutexes, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Furrowlight's program timed out waiting for a booking, sent it again, and got a `409` for the pass it had booked itself seconds before.

A program that may retry a booking now sends an `Idempotency-Key`, a name it makes up for that one request. The service keeps its answer under the key for 24 hours, and a retry with the same key and body gets that answer again and books nothing. The body's SHA-256 hash tells a true retry from a new request that reuses a key:

```go
a := sha256.Sum256([]byte(`{"n":1}`))
b := sha256.Sum256([]byte(`{"n":2}`))
fmt.Println(a == b) // false
```

## Your task

`claim` lets every request run. Under `s.mu`, if `id` was claimed less than 24 hours ago, return `differentRequest(id.Key)` when the body's hash differs, `inProgress(id.Key)` while that request isn't `Done`, and otherwise the `kept` answer. If not, store a new `kept` for `id` and return nil, nil.
