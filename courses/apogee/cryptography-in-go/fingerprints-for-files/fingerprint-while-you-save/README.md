---
slug: fingerprint-while-you-save
title: Fingerprint While You Save
type: code_tests
position: 6
difficulty: 7
topics: crypto, streams
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
The program that delivers files, dlinkd, writes each image once into the folder customers fetch it from. Reading 10.9 GB back from the disk only to fingerprint it would read every byte twice, so Nkem asks for the digest on the way in.

`io.MultiWriter` returns one writer that writes everything to several: give it the file and a hash, and one `io.Copy` saves the bytes and digests them as they pass.

```go
var a, b bytes.Buffer
w := io.MultiWriter(&a, &b)
io.WriteString(w, "PING")
fmt.Println(a.String(), b.String()) // PING PING
```

## Your task

`save` copies `src` to `dst`, then returns `""`. Make it copy into `io.MultiWriter(dst, h)`, with `h` from `sha256.New()`, and return `sha256:` followed by the hex digest, or `""` and the copy's error.
