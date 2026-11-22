---
slug: internal-means-internal
title: internal Means Internal
type: code
position: 4
difficulty: 6
topics: compiler-errors, packages
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:20:22Z
---
Mirela has moved the signing keys to `uplink/internal/keys`, and `main.go` no longer builds. "Good. Keys go through `uplink`, nothing else."

A package inside a folder named `internal` can be imported only from the folder that holds `internal` and the folders below it. The go command enforces this when it builds:

```text
downlink/internal/crc    the package
downlink                 may import it
downlink/assemble        may import it
alert, main              may not
```

## Your task

`main.go` imports `apogee/groundctl/uplink/internal/keys`, and the build fails with `use of internal package ... not allowed`. Import `apogee/groundctl/uplink` instead, and print the ID that `uplink.KeyID` returns.
