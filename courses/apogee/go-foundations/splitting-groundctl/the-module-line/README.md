---
slug: the-module-line
title: The Module Line
type: code
position: 0
difficulty: 5
topics: compiler-errors, modules, packages
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
groundctl, the software behind Apogee's ground stations, has grown into one `main` package that three teams edit at once. Nkem starts splitting it up: the orbit math now lives in its own folder, `orbit/`.

The first line of `go.mod`, `module apogee/groundctl`, sets the module path; from version 2 on, the path ends in the major version, as in `apogee/groundctl/v2`. Each folder is one package, and all its files start with the same `package` line. Import a package by the module path plus its folder, then call it by the name on that line, by convention the folder's name:

```go
import "apogee/groundctl/fleet" // code in fleet/

func main() {
	// fleet/*.go start with package fleet
	fmt.Println(fleet.Size())
}
```

## Your task

`main.go` calls `orbit.PeriodMinutes`, but nothing imports `orbit`, so the build stops with `undefined: orbit`. Import the `orbit` package in `main.go`.

> [!NOTE]
> `go.mod` is read-only: every lesson in this chapter builds the same module. A module that uses other modules lists each one and its version in a `require` line, and `go.sum` beside it keeps a checksum of each one's contents, a number used to check that they have not changed.
