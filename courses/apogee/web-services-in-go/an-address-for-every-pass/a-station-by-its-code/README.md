---
slug: a-station-by-its-code
title: A Station by Its Code
type: code_tests
position: 1
difficulty: 6
topics: http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Yusra wants every station to have an address of its own, so her customer guide can link straight to Salar Quieto at `/stations/SQT`.

A wildcard, a name in braces, matches one whole segment of the path: the text after a slash, up to the next slash or the end. `GET /stations/{code}` matches `/stations/SQT`, but neither `/stations/` nor `/stations/SQT/extra`. In the handler, `r.PathValue` returns the text the wildcard matched:

```go
mux.HandleFunc("GET /dishes/{name}", showDish)

// in showDish, for GET /dishes/chiquita:
name := r.PathValue("name") // "chiquita"
```

## Your task

1. In `routes`, register `GET /stations/{code}` for `s.getStation`.
2. In `getStation`, read the code from the path instead of using `""`. The `findStation` helper does the rest, so the service answers:

```text
GET /stations/SQT   200 and Salar Quieto
GET /stations/ZZZ   404 no station "ZZZ"
GET /stations/sqt   404 no station "sqt"
```

> [!NOTE]
> Codes match exactly: station codes are capitals, and `sqt` is none of them.
