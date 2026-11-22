---
slug: the-log-parcel
title: The Log Parcel
type: code_tests
position: 6
difficulty: 7
topics: files, loops, streams, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:26Z
---
The stick's second parcel holds the controller log since Paz left in July: `lg-log-2027-10-15.tar.gz`, a tar, which strings files one after another, squeezed by gzip. Mirela's note: "Zeros pack a thousandfold. A parcel of one megabyte can unpack into a gigabyte."

Read it in layers: `gzip.NewReader(r)` unsqueezes the bytes, and `tar.NewReader` reads files from them. `Next` returns the next file's header, or `io.EOF` after the last file; reads from the tar reader then give that file's bytes. Read each file through `io.LimitReader` with one byte more than the limit: if that byte arrives, the file is too long.

```go
zr, err := gzip.NewReader(r)
if err != nil {
	return err
}
tr := tar.NewReader(zr)
h, err := tr.Next() // io.EOF: no more files
```

## Your task

Implement `readLogs`: return the text of every file in the parcel by its name, or an error naming the first file longer than `limit` bytes. Return the other errors too, such as the one for a log that isn't gzip'd.
