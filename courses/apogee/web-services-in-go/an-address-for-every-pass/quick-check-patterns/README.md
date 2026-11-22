---
slug: quick-check-patterns
title: Quick Check: Patterns
type: quiz
position: 3
difficulty: 6
topics: http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Hollis's links go out on Friday. Nkem covers most of `routes` with her hand and leaves two lines showing:

```go
mux.HandleFunc("GET /stations", s.listStations)
mux.HandleFunc("GET /stations/{code}",
	s.getStation)
```

"Before anyone clicks a link, tell me what the mux answers."

## Your task

Answer three questions about these two routes alone: what a HEAD request gets, what a POST and an extra segment get, and what `{code}` matches.
