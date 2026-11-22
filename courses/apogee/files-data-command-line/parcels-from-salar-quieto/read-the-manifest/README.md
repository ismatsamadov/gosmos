---
slug: read-the-manifest
title: Read the Manifest
type: code_tests
position: 0
difficulty: 7
topics: structs, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T03:52:00Z
---
Friday, 3,700 meters. A stick has come down from Salar Quieto, 140 meters above the van, with a parcel from its dish controller: one file that packs La Grande's settings together, with `manifest.xml` listing everything inside. Paz wants that list read before anything is unpacked.

The list is XML: plain text in which each value sits between a start tag and an end tag, as in `<dish>LG</dish>`, and a start tag can carry attributes. `xml.Unmarshal` fills a struct through its tags: `xml:"name"` reads a child element, `xml:"name,attr"` an attribute, `xml:"a>b"` every `b` inside an `a`, and `xml:",chardata"` the text of the element the struct stands for.

```go
type Shift struct {
	Crew  int      `xml:"crew,attr"`
	Chief string   `xml:"chief"`
	Notes []string `xml:"notes>note"`
}
```

## Your task

Tag the fields of `Manifest` and `Entry` in `manifest.go` so that `manifest.xml`, shortened here, decodes into them:

```xml
<parcel version="1">
  <dish>LG</dish>
  <made>2027-10-15T09:40:00Z</made>
  <files>
    <file size="153">correction.xml</file>
  </files>
</parcel>
```
