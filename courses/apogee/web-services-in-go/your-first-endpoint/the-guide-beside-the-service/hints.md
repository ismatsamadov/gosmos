`http.FileServerFS(s.Guide)` is a handler that serves the guide's files. Register it with `mux.Handle`, not `HandleFunc`, since it is a handler already, at the pattern `/guide/` with its slash.

The file server looks up the request's whole path, and the guide has no folder called `guide`. Wrap it in `http.StripPrefix("/guide/", ...)`, so `/guide/stations.html` reaches it as `stations.html`.

If every page is 404, the prefix is still on the path, or the pattern lacks its trailing slash; if `/guide` doesn't redirect, the pattern is `/guide` instead of `/guide/`.
