---
slug: louder-without-a-restart
title: Louder Without a Restart
type: code_tests
position: 1
difficulty: 7
topics: http, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:37Z
---
At 08:00 this Monday, Lantern Atoll's uplinkd was restarted: the Monday restart, which Ingrid added to her checklist "until we know why." A restart clears whatever went wrong, so the night shift will turn a station's logging up while it runs, through uplinkd's debug port: a private address on the station's own network, for engineers only.

A `slog.LevelVar` is a level that can change while the program runs, safely from any goroutine. A handler given `&v` as its `Level` checks every record against `v` as it is at that moment. `UnmarshalText` sets `v` from a level's name, such as `debug` or `WARN`, in any letter case; a name it doesn't know is an error, and `v` stays as it was:

```go
var v slog.LevelVar // INFO until set
err := v.UnmarshalText([]byte("Error"))
fmt.Println(err, v.Level()) // <nil> ERROR
err = v.UnmarshalText([]byte("loud"))
fmt.Println(err != nil, v.Level()) // true ERROR
```

## Your task

`levelHandler` serves `PUT /debug/level`; its pattern already answers 405 to any other method. Replace its TODO: set `v` from the body with `UnmarshalText`, and when that fails, answer 400 with the body quoted, `unknown level "loud"`, instead of 204.
