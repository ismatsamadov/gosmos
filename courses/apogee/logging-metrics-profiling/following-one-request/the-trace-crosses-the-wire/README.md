---
slug: the-trace-crosses-the-wire
title: The Trace Crosses the Wire
type: code_tests
position: 5
difficulty: 10
topics: http, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Monday 19 February, the day spans reach every uplinkd. A pass begins at headquarters, when uplinkd there hands Lantern Atoll its lines with `POST /uplink`, and the trace must cross with them.

A context stays in one program, so a span crosses in a header, `Apogee-Trace: <trace id>-<span id>`. Going out, a `RoundTripper` sets it on a clone of each request, never on the caller's own. Coming in, the station puts the header's span, a remote parent with only its IDs, in the request's context and starts its own span as its child:

```go
out := r.Clone(r.Context())
out.Header.Set("Apogee-Pass", id)
fmt.Println(r.Header.Get("Apogee-Pass") == "")
// true: the caller's request is unchanged
```

## Your task

1. In `RoundTrip`, when the request's context carries a span, send a clone with `Apogee-Trace` set from it.
2. Make `continueTrace` put `remoteParent`'s span, if any, in the request's context, start a span called `name` there, serve `next` with the span's context, and end the span.

> [!NOTE]
> `remoteParent`, in `header.go`, refuses anything but 16 and 8 lowercase hex digits around one `-`: such a request gets a trace of its own.
