---
slug: say-it-with-a-status
title: Say It With a Status
type: code_tests
position: 4
difficulty: 6
topics: http, slices
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Yusra's customer guide promises that `/stations?band=X` lists the stations a customer can book on X-band, the radio band that carries big files, and Furrowlight's engineers have already typed `band=x`. They need a clear refusal, not the whole list.

When a handler can't use a request, it says so with a status before writing anything else: `http.Error(w, msg, code)` sets the status and writes `msg` as plain text. Once the body has begun, the status can't change.

```go
if r.URL.Query().Get("units") == "miles" {
	http.Error(w, "we measure in km",
		http.StatusBadRequest)
	return
}
```

## Your task

`listStations` ignores `band`. Read it with `r.URL.Query().Get`: with no band, answer every station; with exactly `UHF`, `S` or `X`, the stations whose `Bands` hold it; with anything else, answer 400 through `http.Error`, quoting the band with `%q`:

```text
unknown band "x"
```

> [!NOTE]
> Salar Quieto and Cape Wren can only listen on UHF, the band of small satellites' beacons, so their `Bands` leave it out.
