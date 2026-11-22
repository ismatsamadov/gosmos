Turn the value into bytes before anything is written: `json.Marshal(v)` returns the JSON and an error. When the error isn't nil, nothing has gone out yet, so `http.Error` can still answer 500.

After a successful marshal the order is the whole lesson: `w.Header().Set` first, then `w.WriteHeader(status)`, then `w.Write` with the JSON and a newline appended.

If the test still sees `text/plain` or a 200, something was written before the header or the status: the first write fixes both, and nothing after it can change them.
