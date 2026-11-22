---
slug: unpack-the-whole-parcel
title: Unpack the Whole Parcel
type: code_tests
position: 8
difficulty: 8
topics: errors, files, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:27Z
---
Sunday night, before Monday's climb. Paz wants La Grande's settings unpacked onto the stick, so that if Thursday's calibration goes wrong, the old correction is there to put back.

First run every check that writes nothing, so a bad parcel leaves the stick untouched. Then write only through the root, every file under its cap, and read the manifest last, through the root, from a copy the cap has already kept small.

```go
if err := checkAll(items); err != nil {
	return err // nothing written yet
}
for _, it := range items {
	if err := write(it); err != nil {
		return err
	}
}
```

## Your task

Implement `unpackParcel(path, dir)` with the supplied helpers:

1. Open the zip, and refuse it with `checkNames` before anything is written.
2. Open `dir` as a root, and unpack every file with `unpackFile`, at most `maxFileBytes` each.
3. Read `manifest.xml` through the root and return what `readManifest` makes of it.
