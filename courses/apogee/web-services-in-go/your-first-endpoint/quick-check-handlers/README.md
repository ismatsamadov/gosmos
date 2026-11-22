---
slug: quick-check-handlers
title: Quick Check: Handlers
type: quiz
position: 2
difficulty: 5
topics: functions, http, methods
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Before the service answers anyone outside Apogee, Nkem wants the handlers understood, not only written. Her whiteboard holds three numbered lines of code and a question for each; `s` is a `*Service`:

```go
var h http.Handler = stationList(nil) // 1
g := http.HandlerFunc(s.listStations) // 2
f := s.listStations                   // 3
```

## Your task

Answer one question for each line: why line 1 compiles, what line 2 does when it runs, and what `f` holds after line 3.
