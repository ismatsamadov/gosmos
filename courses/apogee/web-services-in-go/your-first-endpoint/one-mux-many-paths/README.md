---
slug: one-mux-many-paths
title: One Mux, Many Paths
type: code_tests
position: 5
difficulty: 6
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Ingrid adds the booking service to her monitor, which will ask it every minute whether it is up, at `/health`: "It answers ok, or someone gets woken up."

An `http.ServeMux` hands each request to the handler registered for its path. A pattern without a trailing slash, such as `/stations`, matches that path alone, and any path nobody registered gets the mux's own answer, `404 page not found`:

```go
mux := http.NewServeMux()
mux.HandleFunc("/stations", s.listStations)
rec := httptest.NewRecorder()
mux.ServeHTTP(rec, httptest.NewRequest(
	http.MethodGet, "/statoins", nil))
fmt.Println(rec.Code) // 404
```

## Your task

`routes` registers `/stations` only. Write `health` so that it answers `{"status":"ok"}` with status 200 through `writeJSON`, and register it at `/health`.

> [!WARNING]
> Don't register `/` for the health check: `/` matches every path nobody else claimed, so `/statoins` would answer ok too.
