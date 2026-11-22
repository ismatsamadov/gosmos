`runtime.GC()` takes no arguments and returns once a whole collection is done. Call it after the file is created and before anything is written, so the profile counts only what is still in use.

`pprof.Lookup("heap")` returns the heap profile, and its `WriteTo` takes the file and a debug level: 0 writes the compressed form that `go tool pprof` reads. When `WriteTo` fails, close the file and return its error.

If the test finds an empty file, nothing was written into `f` before it was closed. If it counts 0 garbage collections, `runtime.GC()` is missing from `heapProfile`.
