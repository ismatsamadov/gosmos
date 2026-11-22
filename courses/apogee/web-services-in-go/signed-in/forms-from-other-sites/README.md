---
slug: forms-from-other-sites
title: Forms From Other Sites
type: code_tests
position: 4
difficulty: 9
topics: http, security
created: 2026-11-13T03:49:38Z
modified: 2026-11-13T03:49:38Z
---
Mirela sends you a page she made on her laptop with one button, "See Furrow-A's new image". Behind it is a form that posts to `/account/bookings/bk-4776/cancel`: a browser signed in as Furrowlight would cancel Monday's pass.

A browser adds the service's cookies to every request it sends there, whichever page started it. That is cross-site request forgery: another site's page acting with your session. `SameSite=Lax` keeps the cookie off forms from other sites in current browsers, but not off forms from other pages of the same site, and older browsers ignore it. So the service checks as well: `http.CrossOriginProtection` refuses any request but GET, HEAD and OPTIONS that the browser says came from another origin, in its `Sec-Fetch-Site` or `Origin` header.

```go
p := http.NewCrossOriginProtection()
r := httptest.NewRequest("POST", "/x", nil)
r.Header.Set("Sec-Fetch-Site", "cross-site")
fmt.Println(p.Check(r) != nil) // true: refuse
```

## Your task

`noCrossOrigin` returns `h` unguarded. Return a handler that puts every request through a new `http.CrossOriginProtection` and answers 403 to those it refuses: one of its methods does both.
