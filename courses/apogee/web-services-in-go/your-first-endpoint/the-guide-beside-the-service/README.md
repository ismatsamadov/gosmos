---
slug: the-guide-beside-the-service
title: The Guide Beside the Service
type: code_tests
position: 6
difficulty: 6
topics: files, http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Yusra wants her customer guide served by the booking service itself, at `/guide/`, so every customer reads the same version: the one that matches the service they are calling.

`http.FileServerFS(fsys)` answers a request with the file at its path inside `fsys`, and a folder with its `index.html`. Under a prefix, each path still starts with `/guide/`, so `http.StripPrefix` cuts that off first. A pattern ending in `/` matches every path below it.

```go
files := http.FileServerFS(os.DirFS("photos"))
h := http.StripPrefix("/photos/", files)
// GET /photos/gull.jpg reaches files
// with the path gull.jpg
```

## Your task

`routes` has no guide yet. Register a handler at `/guide/` that serves `s.Guide` with `http.FileServerFS`, behind `http.StripPrefix("/guide/", ...)`.

> [!NOTE]
> A request for `/guide`, without the slash, gets a 307 redirect to `/guide/` from the mux itself.
