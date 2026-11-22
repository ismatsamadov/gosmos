---
slug: unpack-under-a-root
title: Unpack Under a Root
type: code_tests
position: 7
difficulty: 7
topics: errors, files, streams, testing
created: 2026-11-13T03:52:00Z
modified: 2026-11-13T12:10:27Z
---
Sunday. Mirela has one more question: what if the stick's folder already holds a symbolic link, a name that points to another place on the disk? With a link at `log/` pointing outside, `log/2027-10-14.txt` passes `IsLocal`, since its text stays inside, yet the file would land outside.

`os.OpenRoot(dir)` returns an `*os.Root`, a handle on one folder. Its methods, such as `MkdirAll`, `Create` and `ReadFile`, take names inside that folder and refuse any name that would end up outside it: an absolute one, a `..` that climbs out, or a link that leads out.

```go
root, err := os.OpenRoot("stick")
if err != nil {
	return err
}
defer root.Close()
err = root.WriteFile("fog.txt", data, 0o600)
```

## Your task

Complete `unpackFile`: make the folder `f.Name` sits in with `root.MkdirAll` and mode `0o700`, create the file with `root.Create`, and copy `src` into it through `io.LimitReader`, refusing a file longer than `limit` as you did for the logs. Return the first error, `Close`'s included.
