---
slug: hollis-signs-in
title: Hollis Signs In
type: code_tests
position: 8
difficulty: 9
topics: http, security
created: 2026-11-13T03:49:39Z
modified: 2026-11-13T12:18:59Z
---
Friday, 18:00 UTC. Hollis is waiting to sign in, and Yusra wants the page in the beta tonight. Every guard is written; now each one goes where it belongs.

`site` splits the addresses with a mux of its own: `/account/` to the account page, everything else to the API, whose own mux still answers its 404s and 405s. The page's callers show a cookie, which the browser sends by itself, so the page needs `noCrossOrigin`, and Mirela's headers too. Like the API, it needs `noStore`: a kept copy of a signed-in page is someone's bookings. The API's callers show a token, which no other site's form can carry, so it needs only `cors`.

```go
// the first one listed runs first
h := chain(mux, requestID, logRequests(l, now))
```

## Your task

`site` serves both halves unguarded. Put the account page's routes behind `securityHeaders`, `noStore` and `noCrossOrigin`, in an order that gives its refusals the headers too, and the API behind `cors(origin)` alone.
