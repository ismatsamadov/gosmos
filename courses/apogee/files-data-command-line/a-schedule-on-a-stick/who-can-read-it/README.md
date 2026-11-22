---
slug: who-can-read-it
title: Who Can Read It
type: code_tests
position: 5
difficulty: 3
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-04T17:28:41Z
---
Mirela checks the van's reports before departure. A write can replace a file's contents while keeping its old permissions, so a successful save alone does not prove who can read it.

`os.Stat` returns information about a path. `info.Mode().Perm()` selects its nine permission bits: owner, group and others. A bitwise `&` keeps only the bits selected by a mask:

```go
mode := info.Mode().Perm()
groupBits := mode & 0070
fmt.Printf("group permissions: %o\n", groupBits)
```

## Your task

Implement `ownerOnly`: inspect the path and return whether its group and other bits, selected by `0077`, are all zero. Return a Stat error rather than treating a missing file as private. This checks permission bits, not every possible operating-system access rule.
