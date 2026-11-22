Write the method in place of the TODO line, on `APIKey` itself rather than a pointer, named `LogValue`, with no arguments, returning an `slog.Value`. Add `log/slog` and `strings` to the imports.

`string(k)` gives the key's text, and `strings.Cut` splits it at the first dot and reports whether there was one. Make the value with `slog.StringValue`: the customer ID followed by `.[redacted]`, or `[redacted]` alone.

If the test says `APIKey` has no `LogValue` method, check the name, the signature and the receiver: with a pointer receiver, a key logged by value would never find the method.
