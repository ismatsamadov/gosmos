---
slug: one-trusted-origin
title: One Trusted Origin
type: code_tests
position: 5
difficulty: 9
topics: http, security
created: 2026-11-13T03:49:39Z
modified: 2026-11-13T12:18:56Z
---
Yusra's TL;DR: "Furrowlight's planning page wants Monday's passes straight from the API." The planning page, where Furrowlight's engineers choose which passes to book, lives at `https://plan.furrowlight.example`. Its script calls `GET /passes` with a token, and the browser hides every answer from it.

A browser lets a page's script read an answer from another origin only when the answer says so in `Access-Control-Allow-Origin`: that is CORS. A request carrying `Authorization` first gets a preflight, an `OPTIONS` that asks which method and headers the service allows:

```text
OPTIONS /passes HTTP/1.1
Origin: https://plan.furrowlight.example
Access-Control-Request-Method: GET
Access-Control-Request-Headers: authorization
```

## Your task

`cors` passes requests on untouched. Add `Vary: Origin` to every answer, since each now depends on `Origin`. When `Origin` is exactly `origin`, set `Access-Control-Allow-Origin` to it, and answer an `OPTIONS` yourself with 204 and these two headers, without calling `next`:

1. `Access-Control-Allow-Methods: GET, POST, DELETE`
2. `Access-Control-Allow-Headers: Authorization, Content-Type, Idempotency-Key`

Everything else goes on to `next`.

> [!WARNING]
> Never `*`, and never whatever `Origin` arrives: either would let any site's pages call the API.
