---
slug: four-headers-for-the-page
title: Four Headers for the Page
type: code_tests
position: 7
difficulty: 8
topics: http, security
created: 2026-11-13T03:49:39Z
modified: 2026-11-13T12:18:58Z
---
Mirela's last list for the page has four lines. Another site could show the page inside its own, under a fake button, so that a click lands on Cancel. A browser could guess an answer's type and run text as a script. The page's address could follow a link to another site. And a visit that starts over plain HTTP, from an old link or a typed address, could be read on the way.

A browser guards against each one only when the answer asks it to, in a header. Headers that must be on every answer are middleware's work, like `noStore`'s:

```go
// no page of ours may use the camera
w.Header().Set("Permissions-Policy", "camera=()")
```

## Your task

`securityHeaders` passes requests on untouched. Before calling `next`, set:

1. `Content-Security-Policy` to `default-src 'self'; frame-ancestors 'none'`: load nothing from elsewhere, and appear in nobody's frame.
2. `Strict-Transport-Security: max-age=31536000`: for a year, this browser reaches the service only over HTTPS.
3. `X-Content-Type-Options: nosniff` and `Referrer-Policy: same-origin`, so the browser takes each type as given and tells no other site where its visitor came from.
