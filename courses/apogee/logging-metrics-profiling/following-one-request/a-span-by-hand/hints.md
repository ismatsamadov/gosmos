Keep the starter's new span in a variable and fill it in before returning it. `spanFrom(ctx)` returns the parent and true when there is one: copy its `TraceID` into the span, and its `SpanID` into the span's `ParentID`.

Lock `t.mu` before the first draw and unlock it after the second: only a root draws a trace ID, with `Uint64`, and then every span draws its own ID with `Uint32`, in that order. `passFrom` returns the pass ID and a bool you can ignore here.

If the first trace isn't `c4f5a58656eef510`, the draws come in another order, or there are more of them. If `spanFrom` finds nothing in the context `Start` returned, return the one `withSpan(ctx, s)` makes, not `ctx` itself.
