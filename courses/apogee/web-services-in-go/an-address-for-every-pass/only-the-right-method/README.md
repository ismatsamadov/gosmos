---
slug: only-the-right-method
title: Only the Right Method
type: code_tests
position: 0
difficulty: 6
topics: http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Furrowlight's engineers started trying the booking service on Monday, and their test program sent `POST /stations` by mistake. The service answered with the station list, as if the POST had changed something.

A pattern can start with a method and a space. `GET /stations` matches GET requests, and HEAD ones, which ask for the same answer without its body. For any other method on that path, the mux answers `405 Method Not Allowed` itself, with an `Allow` header naming the methods that work:

```go
mux.HandleFunc("DELETE /cache", clearCache)
// GET /cache: 405, Allow: DELETE
```

## Your task

`routes` registers `/health`, `/stations` and `/guide/` for every method. Put `GET` in front of each pattern, so that a POST to any of them gets:

```text
405 Method Not Allowed
Allow: GET, HEAD
```

> [!NOTE]
> A path that no pattern matches still gets `404 page not found`. A 405 says the path exists.
