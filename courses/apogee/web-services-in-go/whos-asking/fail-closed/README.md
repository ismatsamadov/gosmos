---
slug: fail-closed
title: Fail Closed
type: code_tests
position: 7
difficulty: 9
topics: context, http, security
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela signs off every review the same way: "LGTM once it fails closed." A guard that fails closed refuses whenever anything is missing or wrong; one that fails open lets a request through when its check goes wrong.

`routes` puts `authenticate` around every private route: it either answers 401 and returns, or calls `next` with the customer in the request's context. The go/no-go poll's wind check fails closed too: no reading is a NO-GO, never a GO.

```go
func windGo(windKmh float64, err error) bool {
	if err != nil {
		return false // no reading: NO-GO
	}
	return windKmh < 72
}
```

## Your task

`authenticate` lets every request through. Read the `Bearer` credential, check it with `verifyToken` at `now()`, and find its customer in `customers`, where an unknown one is `unauthorized("invalid token")`. On any failure, set `WWW-Authenticate: Bearer`, answer through `writeError` and return. Otherwise call `next` with `r.WithContext(withCustomer(r.Context(), c))`.

> [!WARNING]
> Writing the 401 and calling `next` anyway still runs the handler; the test's handler counts its calls.
