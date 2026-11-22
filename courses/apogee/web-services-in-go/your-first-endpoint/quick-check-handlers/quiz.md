## Why does line 1 compile?

- [ ] Every slice type is a handler
  > A slice type gets no methods by itself; `stationList` is a handler only because of the method written for it.
- [x] `stationList` has a `ServeHTTP` method with a handler's two parameters
  > That one method is all `http.Handler` asks for, so any type that has it is a handler.
- [ ] A `nil` value fits into any interface
  > A typed `nil` still has to be a type that satisfies the interface; its method is what makes it one.

## What does line 2 do when it runs?

- [x] It converts the method to a type that has `ServeHTTP`; nothing else runs
  > `http.HandlerFunc(...)` is a conversion, and `listStations` runs only when a request arrives.
- [ ] It calls `listStations` once and keeps what it wrote
  > There is no request to call it with: a conversion runs nothing.
- [ ] It registers `listStations` with the default mux
  > Registering takes `Handle` or `HandleFunc`; a conversion registers nothing.

## After line 3, what does `f` hold?

- [ ] A copy of the stations `s` holds
  > Nothing is copied from the service: `f` is a function, not a list.
- [ ] What `listStations` wrote when line 3 ran
  > Without parentheses there is no call, so nothing was written.
- [x] A method value: a function that runs `listStations` on this very `s`
  > A method value remembers its receiver, so `f(w, r)` is the same call as `s.listStations(w, r)`.
