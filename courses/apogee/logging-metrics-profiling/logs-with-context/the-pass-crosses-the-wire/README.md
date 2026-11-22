---
slug: the-pass-crosses-the-wire
title: The Pass Crosses the Wire
type: code_tests
position: 4
difficulty: 8
topics: context, http, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Thursday. Headquarters' uplinkd hands each pass's commands to Lantern Atoll with `POST /uplink`, and the pass ID stops at the door: a context lives inside one program.

So the ID crosses as a header, `Apogee-Pass`. On the way out, an `http.RoundTripper` wrapped around the client's transport sets it, on a clone of the request, since a `RoundTrip` must never change the request it is given. On the way in, middleware reads the header and puts the ID into the request's context, for the handlers after it:

```go
// inside a RoundTrip method:
r = r.Clone(r.Context()) // a copy to change
r.Header.Set("User-Agent", "uplinkd")
return t.base.RoundTrip(r)
```

## Your task

1. In `RoundTrip`, when `r`'s context carries a pass, send a clone of `r` with `Apogee-Pass` set to its ID through `t.base`.
2. In `acceptPass`, when the request's `Apogee-Pass` isn't empty, call `next` with the request carrying that ID in its context.

> [!NOTE]
> A missing or empty header means no pass: the station's log lines then carry no `pass_id`.
