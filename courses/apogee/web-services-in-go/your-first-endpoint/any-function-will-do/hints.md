`listStations` already has a handler's two parameters, an `http.ResponseWriter` and a `*http.Request`. What it lacks is a `ServeHTTP` method, which `http.HandlerFunc` supplies.

Write `s.listStations` without parentheses, so you pass the method instead of calling it, and convert it with `http.HandlerFunc(...)`, the way you would convert a number with `float64(...)`.

If both services answer with all five stations, the handler isn't using its own `s`: build it from `s.listStations`, not from a new `Service`.
