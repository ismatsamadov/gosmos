## When does the code above set `Content-Language`?

- [ ] Once, when `english` wraps the mux at start
  > Wrapping only builds the function `f` and returns it; nothing inside `f` runs until a request arrives.
- [x] Each time a request arrives, before the mux answers
  > The service calls `f` for every request, and `f` sets the header before it hands the request to `next`.
- [ ] Each time a request arrives, after the mux has answered
  > `f` sets the header before calling `next`; after the answer, it would come too late.

## A handler calls `r2 := r.WithContext(ctx)`. What happens to `r`?

- [x] Nothing: `r2` is a new request carrying `ctx`, and `r` keeps its old context
  > That is why the next handler must get `r2`; handed `r`, it never sees the new value.
- [ ] `r` now carries `ctx` too, like `r2`
  > `WithContext` never changes the request it's called on; only the copy it returns carries `ctx`.
- [ ] `r` is closed, and only `r2` may be used from now on
  > Both requests stay usable; they differ only in the context each one carries.

## Why is the request ID stored under `requestIDKey{}` and not under the string `"request_id"`?

- [ ] Context keys have to be structs
  > Any comparable value can be a key, a string included; that is exactly the danger.
- [ ] A struct key makes looking up the value faster
  > Speed isn't the reason; what matters is who else can make the same key.
- [x] No other package can make a value of that type, so no one else's value can land on the same key
  > Any package could store its own value under the string `"request_id"`, and the two would clash.
