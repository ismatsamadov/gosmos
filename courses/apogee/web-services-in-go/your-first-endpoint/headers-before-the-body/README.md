---
slug: headers-before-the-body
title: Headers Before the Body
type: code_tests
position: 3
difficulty: 6
topics: http, json
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Cirrine's engineers send Nkem their new client to try before Friday, and it refuses the station list: the service sets `Content-Type` after writing the body, so the list goes out labeled `text/plain`.

An answer leaves in order: the status line, the headers, then the body. The first `Write` sends the status and headers as they are at that moment, 200 unless `WriteHeader` chose another, so set the headers, then call `WriteHeader`, then write. Marshal before all of that, so a failure can still change the status.

```go
w.Header().Set("Retry-After", "60")
w.WriteHeader(http.StatusServiceUnavailable)
fmt.Fprintln(w, "mast silent")
```

## Your task

`writeJSON` writes the body first, so its header and status are lost. Marshal `v` first with `json.Marshal`; if that fails, answer 500 with `http.Error` and the text `internal error`. Otherwise set `Content-Type: application/json`, call `WriteHeader(status)`, then write the JSON and a newline.

> [!NOTE]
> `listStations` already answers through `writeJSON`, so this one function fixes the station list.
