---
slug: quick-check-decoders
title: Quick Check: Decoders
type: quiz
position: 3
difficulty: 5
topics: errors, streams
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Mirela reviews gctl's input boundary before the van climbs on, and starts the way her reviews always start: "Show me the file it refuses." Her questions are about this decoder, reading a settings file `r`:

```go
dec := json.NewDecoder(r)
dec.DisallowUnknownFields()
var s Settings
err := dec.Decode(&s)
_, next := dec.Token()
```

## Your task

Answer Mirela's three questions. Each describes a settings file that is slightly different and asks what this code makes of it.
