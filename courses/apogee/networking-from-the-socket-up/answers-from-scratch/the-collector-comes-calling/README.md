---
slug: the-collector-comes-calling
title: The Collector Comes Calling
type: code_tests
position: 8
difficulty: 10
topics: http, networking, switch
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T19:06:41Z
---
Tuesday 20 June, midsummer. Ingrid watches for Furrow-A's 06:32 pass, the first Furrowlight pass through the station server: its image must reach Furrowlight within a minute of the satellite setting, LOS.

With every route in one place, the server can tell a path it doesn't have, 404 `no such path`, from a method that path doesn't take, 405 `method not allowed` with an `Allow` field listing the methods it does. A recording's path needs a valid pass ID; any other target is unknown, a query included:

```text
/recordings              GET, HEAD
/recordings/<pass ID>    GET, HEAD, PUT
/feed                    GET
```

## Your task

`route` sends every request to the list. Send each path's methods to `list`, `recording`, `put` or `feedAnswer()`. With the given `refuse`, answer a known path's other methods 405 `method not allowed` with its `Allow` field, and any other target 404 `no such path`.

> [!NOTE]
> Other methods never reach `route`: `readRequest` answers them 501.
