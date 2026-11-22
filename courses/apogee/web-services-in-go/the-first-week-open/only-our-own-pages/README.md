---
slug: only-our-own-pages
title: Only Our Own Pages
type: code_tests
position: 6
difficulty: 9
topics: http, security, strings
created: 2026-11-13T06:35:19Z
modified: 2026-11-13T16:15:49Z
---
Friday 4 February, 08:30. Mirela forwards an email she wrote herself: "Your API key expires today." Its link starts with the booking service's own address and ends on a page that asks for the key:

```text
/go?next=https://apogee-bookings.example/renew
```

`/go` is how Yusra's emails send customers to the guide: it logs the click and redirects to `next`. A redirect that goes wherever a link says lends the service's name to every such email, so `next` must be a path on this site, one that starts with a single `/`. Browsers read `//host` and `/\host` as another site, and drop tabs and line breaks before they read an address at all; `url.Parse` refuses those.

```go
u, _ := url.Parse("//elsewhere.example/x")
fmt.Println(u.Host) // elsewhere.example
```

## Your task

`localPath` accepts every `next`. Return true only when `next` parses with `url.Parse`, starts with `/` but not with `//`, and holds no `\`.
