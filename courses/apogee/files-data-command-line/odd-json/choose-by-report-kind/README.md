---
slug: choose-by-report-kind
title: Choose by Report Kind
type: code_tests
position: 5
difficulty: 5
topics: pointers, structs, switch, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T17:07:41Z
---
Paz reads Furrowlight's reports for signs of La Grande's drift: a dish that points slightly off loses some of the numbered pieces an image comes down in, its chunks. Image reports count them per file, station reports per station and day.

Each kind needs its own type: decode the envelope, let its kind choose the destination, then decode the raw payload into it. A kind the code doesn't know is an error, never a guess.

```go
var dest any
switch msg.From {
case "YST":
	dest = &Wind{}
case "LAN":
	dest = &Rain{}
}
err := json.Unmarshal(msg.Body, dest)
```

## Your task

Implement `decodeReport` with the read-only `decodeEnvelope`: return an `*ImageReport` for the kind `"image"` and a `*StationReport` for `"station"`, each decoded from the payload. Return an error naming any other kind, and any decoding error.
