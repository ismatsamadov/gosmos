Declare an empty struct type such as `requestIDKey`, and use its value, the type's name followed by `{}`, as the key both when you store the ID and when you look it up.

`r.WithContext` doesn't change `r`; it returns a new request. Hand that new request to `next.ServeHTTP`, or the handler inside never sees the ID.

In `requestIDFrom`, `ctx.Value` returns an `any`. Take the string out with a type assertion in its two-value form, which gives an empty string instead of a panic when nothing is stored under your key.
