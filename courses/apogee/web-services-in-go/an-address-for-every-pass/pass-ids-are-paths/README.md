---
slug: pass-ids-are-paths
title: Pass IDs Are Paths
type: code_tests
position: 2
difficulty: 6
topics: http, strings
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Hollis wants a link to each pass in his emails. Every pass already has a name made of its satellite, its station and the number of the satellite's orbit, called its pass ID, and one such as `furrow-a/SQT/15864` reads like a path.

One wildcard matches one segment, so `GET /passes/{id}` never matches an ID with slashes in it. Give each segment a wildcard of its own, and read each with `PathValue`:

```go
// GET /logs/{station}/{day}
key := r.PathValue("station") + "@" +
	r.PathValue("day")
// /logs/YST/2027-11-24 gives "YST@2027-11-24"
```

## Your task

`getPass` is registered at `GET /passes/{id}`, which `/passes/furrow-a/SQT/15864` never matches. Register it at `GET /passes/{sat}/{station}/{orbit}` instead, and make `getPass` join the three values with `/` into the ID it looks for in `s.Catalog`:

```text
GET /passes/furrow-a/SQT/15864
200 and the pass
GET /passes/furrow-a/SQT/15865
404 no pass furrow-a/SQT/15865
```

> [!TIP]
> `PassView` embeds `Pass`, so its JSON puts `id` beside `sat_id`, `aos` and the pass's other fields.
