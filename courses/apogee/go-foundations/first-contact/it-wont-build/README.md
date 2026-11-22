---
slug: it-wont-build
title: It Won't Build
type: code
position: 5
difficulty: 2
topics: compiler-errors, programs
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
AOS is in six minutes, and the banner tool on Ingrid's console won't build.

Before anything runs, the Go compiler checks the whole program. If it finds a problem, nothing runs. Instead, you get an error that starts with where the problem is, as file, line and column, and then says what's wrong:

```text
./main.go:4:8: "time" imported and not used
```

In Go, importing a package you never use is an error, not a warning. Several imports can share one `import`, in parentheses, one per line.

## Your task

Read the starter's error and remove what it names, so the banner builds and prints:

```text
Banner: Pipit-1 over Port Ellery, AOS 04:12 UTC
```
