---
slug: the-changing-report
title: The Changing Report
type: code_tests
position: 8
difficulty: 6
topics: errors, streams, structs, switch, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T17:07:41Z
---
Furrowlight switched its delivery reports to version 3 on Thursday, while the van had no signal, and nobody on board knows what changed. Yusra needs every report gctl can't read named by its exact ID, to ask for it again in version 2.

Read a changing format in order: decode the envelope strictly, check its version, then decode the payload for its kind. Version 1 reports, from July, have no kind: all of them are image reports.

```go
switch m.Version {
case 1:
	m.Unit = "km/h" // version 1 had no unit
case 2:
default:
	return errors.New("unknown version")
}
```

## Your task

Implement `normalizeReport` with the read-only `decodeOne` and `decodePayload`. Return a `Report` for versions 1 and 2, giving version 1 the kind `"image"`. For any other version, return an error such as `report 9007199254741007: unsupported version 3`, and wrap payload errors with the ID the same way.
