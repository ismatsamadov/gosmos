---
slug: quick-check-xml
title: Quick Check: XML
type: quiz
position: 3
difficulty: 7
topics: structs
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:24Z
---
Late on Friday the van finds a signal, and Nkem reviews the XML code from Port Ellery, one question at a time. Each is about this type, and about `s := Stow{Station: "CWR", WindKmh: 74}`:

```go
type Stow struct {
	XMLName xml.Name `xml:"stow"`
	Station string   `xml:"station,attr"`
	WindKmh int      `xml:"wind_kmh"`
}
```

## Your task

Answer Nkem's three questions: which XML fills `Station`, what decoding does with an element it doesn't expect and one it can't find, and what `xml.MarshalIndent` writes first.
