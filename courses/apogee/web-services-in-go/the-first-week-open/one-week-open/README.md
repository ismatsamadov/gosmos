---
slug: one-week-open
title: One Week Open
type: code_tests
position: 8
difficulty: 10
topics: http, security
created: 2026-11-13T06:35:19Z
modified: 2026-11-13T16:15:50Z
---
Friday, 16:00. The week's work goes out in tonight's deploy, and Nkem wants it served as one handler: the new routes beside the API, inside the same middleware, each behind the guard it needs.

A route that acts for a customer takes a token and counts against that customer's limit, as the API's routes do through `s.private`:

```go
mux.Handle("GET /passes",
	s.private(s.listPasses))
```

A download link is its own credential, so `/files/{id}` takes no token, and neither does `/go`. Every booking counts against its customer's plan, and no webhook dials inside a network.

## Your task

`firstWeek` serves the new routes unguarded. Fix it:

1. Put the upload, download and link routes behind `s.private`; leave `/files/{id}` and `/go` open.
2. Make `s.Store` a `metered` around the store it has, counting against `q`.
3. Give the notifier's `Sender` the client `webhookClient()` returns.
