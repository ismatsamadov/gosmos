---
slug: open-the-parcel
title: Open the Parcel
type: code_tests
position: 4
difficulty: 7
topics: errors, files, loops, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:25Z
---
Saturday morning. So far gctl has read a copy of `manifest.xml` that Paz took out of the parcel by hand. Now it must open the parcel itself, `lg-settings-2027-10-15.zip`: a zip file, the usual way to pack many files into one, each squeezed smaller.

`zip.OpenReader(path)` opens a zip file and reads its directory, the list of every file it holds; `Close` lets the zip go. That list is its `File` field: each `*zip.File` has a `Name`, and its `Open` method returns a reader that unpacks the file's bytes as you read them. Close that reader too.

```go
zr, err := zip.OpenReader("photos.zip")
if err != nil {
	return err
}
defer zr.Close()
for _, f := range zr.File {
	fmt.Println(f.Name, f.UncompressedSize64)
}
```

## Your task

Implement `parcelManifest(path)`: open the zip at that path, find the file named `manifest.xml`, read it whole and return what the supplied `readManifest` makes of it. A parcel with no `manifest.xml` is an error.
