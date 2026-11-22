---
slug: the-board-as-data
title: The Board as Data
type: code_tests
position: 0
difficulty: 7
topics: http, json, mutexes
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
Tuesday. Reloading the whole page every fifteen seconds makes the kitchen screen flicker, and Ingrid wants the page to fetch only the data.

`Refresh` already turns each view into JSON and keeps the bytes in `b.json`, so the data has one source: serve those bytes as they are, read under `b.mu` as `current` reads the view. `Cache-Control: no-cache` lets a screen keep a copy but tells it to ask before using it, which the next lesson builds on. `script.html` now fetches the data every 15 seconds and redraws the rows, the panels and the on-call line.

## Your task

`GET /board.json` answers 404.

1. In `Handler`, register `GET /board.json` with `b.serveJSON`.
2. In `json.go`, write `serveJSON`: set `Content-Type: application/json` and `Cache-Control: no-cache`, then write `b.json`, which starts:

```text
{"passes":[{"id":"bk-4988","sat_id":"pipit-1",
```
