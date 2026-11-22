---
slug: what-the-decoder-skips
title: What the Decoder Skips
type: code_tests
position: 1
difficulty: 7
topics: errors, structs, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:24Z
---
The manifest also holds `<firmware>5.3</firmware>`, the version of the controller's own software, and gctl read straight past it. Nkem, at the evening signal: "Good. Now what does it read when `<dish>` is missing?"

`xml.Unmarshal` skips every element and attribute that no field asks for, and a field whose element is missing keeps its zero value. Neither is an error, and unlike `json.Decoder` it has no switch that refuses unknown names, so the checks for what must be there come after decoding:

```go
var s Shift
if err := xml.Unmarshal(data, &s); err != nil {
	return Shift{}, err
}
if s.Chief == "" {
	return Shift{}, errors.New("no <chief>")
}
```

## Your task

Complete `readManifest` in `main.go`: after decoding, return an error for a manifest with no `<dish>`, no `<made>` or no `<file>`, naming what is missing, such as `manifest: no <dish>`. A `time.Time` that was never set reports true from `IsZero`.
