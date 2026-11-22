---
slug: names-that-stay-inside
title: Names That Stay Inside
type: code_tests
position: 5
difficulty: 7
topics: errors, files, testing
created: 2026-11-13T03:51:59Z
modified: 2026-11-13T12:10:25Z
---
Mirela reads the new code from Port Ellery and sends test parcels of her own. "A name in a zip is whatever its packer typed. Show me the parcel gctl refuses."

Joined to the folder it is unpacked into, a name can lead anywhere: `/etc/hosts` is absolute, and `../gctl` climbs out of the folder. `filepath.IsLocal(name)` reports whether a name stays inside any folder it is joined to: it is not empty, not absolute, and never climbs above where it started. It reads only the text, so every name can be checked before anything is written.

```go
names := []string{"a/x", "a/../x", "../x"}
for _, n := range names {
	fmt.Println(n, filepath.IsLocal(n))
}
// a/x true, a/../x true, ../x false
```

## Your task

Implement `checkNames`: return an error for the first file whose name `filepath.IsLocal` refuses, quoting that name, such as `unsafe name "../gctl"`. Return nil when every name stays inside.
