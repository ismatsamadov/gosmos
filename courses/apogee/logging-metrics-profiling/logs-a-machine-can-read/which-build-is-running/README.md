---
slug: which-build-is-running
title: Which Build Is Running
type: code_tests
position: 7
difficulty: 6
topics: logging, modules
created: 2026-11-13T11:32:35Z
modified: 2026-11-13T16:25:47Z
---
Thursday, 16:00. Ingrid's checklist for Friday's deploy at Lantern Atoll has one new step: "14:01. LAN uplinkd: first line names the new build. Confirm." Every program's first record must say which build it is.

`debug.ReadBuildInfo`, from `runtime/debug`, returns what the go command wrote into a program as it built it, and `ok`, false only for a program built without modules. `GoVersion` is the Go release that built it; `Main` is its module, with the `Path` from `go.mod` and a `Version` from version control, such as `v1.8.0`, or `(devel)` outside it, as in your lesson folder:

```go
bi, ok := debug.ReadBuildInfo()
if ok {
	fmt.Println(bi.GoVersion, bi.Main.Path)
}
// go1.27.3 lesson
```

## Your task

Add three attributes from the build info to the record `logStart` writes: `go_version`, `module` and `module_version`, from `GoVersion`, `Main.Path` and `Main.Version`. When `ok` is false, keep the message alone.
