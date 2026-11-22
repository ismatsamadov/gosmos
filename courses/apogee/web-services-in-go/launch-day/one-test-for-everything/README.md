---
slug: one-test-for-everything
title: One Test for Everything
type: code_tests
position: 7
difficulty: 10
topics: http, testing
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Monday 31 January, the full rehearsal. Nkem wants one test that would have caught every bug of the season, run against the whole service the way a customer's program meets it: over HTTP.

`client.go` has helpers: `request` sends one request, with a token, a body and headers if you give them; `token` trades an API key for a token; and an answer's `errorCode` reads the code of a JSON error, `""` for anything else:

```go
a := request(t, base, "GET", "/stations/ZZZ",
	"", nil)
if a.Status != 404 {
	t.Errorf("GET /stations/ZZZ = %d, want 404",
		a.Status)
}
```

## Your task

In `check.go`, write `checkService` against `base`:

1. `/health` answers 200, and Furrowlight and Brinewatch get tokens.
2. Furrowlight lists its passes and books the first one twice with one `Idempotency-Key`: the second answer is the first, byte for byte.
3. Brinewatch's GET of that booking is 404 `not_found`, and within 11 more requests it gets 429 `rate_limited`.
