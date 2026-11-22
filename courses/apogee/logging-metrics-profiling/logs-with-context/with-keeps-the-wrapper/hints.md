`h.Handler.WithAttrs(attrs)` already returns the inner handler with the attributes added; all that is missing is your wrapper around it.

Return `passHandler{...}` with that call inside the braces, and do the same in `WithGroup` with `h.Handler.WithGroup(name)`.

If the run stops with `fatal error: stack overflow`, a method calls itself: inside `WithAttrs`, `h.WithAttrs` is your own method again, while `h.Handler.WithAttrs` is the inner handler's.
