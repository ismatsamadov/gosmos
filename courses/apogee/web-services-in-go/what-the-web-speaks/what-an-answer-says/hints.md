`standIn()` returns a handler, and a handler answers through its `ServeHTTP` method. Call that method on what `standIn()` returns.

`ServeHTTP` takes the writer first and the request second: here the recorder, then the request. The call takes the place of the `_ = req` line.

If the output still shows only `HTTP/1.1 200 OK`, the stand-in never ran before `printResponse` read the recorder: make the call above `printResponse`, not below it.
