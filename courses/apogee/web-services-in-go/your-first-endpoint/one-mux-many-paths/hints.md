`health` has a handler's two parameters already; its body needs one call to `writeJSON` with `http.StatusOK` and a value that encodes as `{"status":"ok"}`, such as a `map[string]string`.

Then register it in `routes` the way `/stations` is registered, with `mux.HandleFunc`, the path `/health` and the method value `s.health`.

If `/statoins` or `/` answers 200, something was registered at `/`, which matches every path: register the exact path `/health` instead.
