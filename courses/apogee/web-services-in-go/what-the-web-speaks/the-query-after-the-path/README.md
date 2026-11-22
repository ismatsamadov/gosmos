---
slug: the-query-after-the-path
title: The Query After the Path
type: code_tests
position: 2
difficulty: 4
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Hollis, who runs operations at Furrowlight, writes in capitals: "ONLY SALAR QUIETO. ONLY FURROW-A. FROM 09:00 OUR TIME." Furrowlight's program will send those wishes as the query, the part after a `?`: `key=value` pairs such as `station=SQT`, joined by `&`.

`url.Values` collects the pairs, and its `Encode` method writes them out with the keys sorted and every value escaped, so a `:` or a `+` arrives as sent:

```go
q := url.Values{}
q.Set("note", "gulls + fog")
q.Set("band", "S")
fmt.Println(q.Encode())
// band=S&note=gulls+%2B+fog
```

## Your task

`passesPath` glues all three values onto `/passes?` by hand. Rewrite it with `url.Values`: set `sat`, `station` and `from` only when they aren't empty, and return `/passes` alone, with no `?`, when all three are empty.

> [!WARNING]
> In a query, `+` stands for a space, as the example's output shows. Glued in by hand, Hollis's time `2027-12-06T09:00:00+01:00` reaches the service with a space where its `+` was.
