---
slug: the-delivery-manifest
title: The Delivery Manifest
type: code_tests
position: 7
difficulty: 7
topics: crypto, files, streams
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Friday 21 July. Cirrine's first delivery with fingerprints lands this afternoon, and Yusra wants every file's state in one report before Cirrine's archive checks them.

Beside the files lies a list of them, the manifest: per line a digest, a size in bytes and a file's name. Check the size first, since it costs nothing, then hash the file with its digest's own hash. `fs.Stat` gives the size:

```go
fsys := fstest.MapFS{
	"echo.txt": {Data: []byte("PING")},
}
info, err := fs.Stat(fsys, "echo.txt")
fmt.Println(info.Size(), err) // 4 <nil>
```

## Your task

`checkManifest` returns nil. Make it report each manifest line in order, as one of these, with `parseDigest` and `fingerprintAs` from the read-only `digest.go`:

```text
ok <file>
wrong size <file>
changed <file>
missing <file>
```

A line that isn't three fields, with a digest `parseDigest` accepts and a whole number of bytes, stops the check with an error naming its number:

```text
manifest line 3: want <digest> <bytes> <file>
```
