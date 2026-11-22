`context.WithCancel` adds no deadline: the context it makes ends only when it, or its parent, is canceled. What should end a pass instead?

The `context` package has a constructor that takes a parent and a `time.Time`, and returns the same two values as `passContext`.
