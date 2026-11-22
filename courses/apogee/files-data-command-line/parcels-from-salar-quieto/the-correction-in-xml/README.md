---
slug: the-correction-in-xml
title: The Correction in XML
type: code
position: 2
difficulty: 7
topics: printing, structs
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T03:52:00Z
---
On Thursday, La Grande's controller takes its new correction, the small angles it adds to every pointing of the dish, only from a file shaped like the parcel's `correction.xml`. Paz wants proof today: gctl writes the current correction, and it must match that file line for line.

`xml.Marshal` names the root element after the Go type and writes everything on one line. An `XMLName xml.Name` field with a tag gives the root its name, and `xml.MarshalIndent(v, prefix, indent)` puts each element on a line of its own. Neither writes the `<?xml …?>` line: `xml.Header` holds it, newline included.

```go
type Gust struct {
	XMLName xml.Name `xml:"gust"`
	Kmh     int      `xml:"kmh"`
}

g := Gust{Kmh: 81}
out, err := xml.MarshalIndent(g, "", "\t")
```

## Your task

1. Give `Correction` an `XMLName` field, so the root element is `correction`.
2. Write it with `xml.MarshalIndent`, indenting each level by two spaces.
3. Print `xml.Header` before it.
