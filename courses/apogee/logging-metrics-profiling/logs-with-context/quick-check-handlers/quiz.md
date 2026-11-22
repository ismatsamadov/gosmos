## Which context does `Handle` get for the third line's record?

- [x] A background context, so the record carries no pass ID
  > `Info` has no context to pass, so the logger hands `Handle` `context.Background()`, which holds no pass.
- [ ] `ctx`, the context the last `InfoContext` call used
  > A logger keeps no context between calls; each record gets the context of its own call.
- [ ] None at all, so `Handle` must check for nil before reading it
  > `Handle` always gets a context; for `Info` it is a background one, never nil.

## If `WithAttrs` returned `h.Handler.WithAttrs(attrs)`, what would a logger made with `With` lose?

- [ ] The attributes from `With`, since the inner handler never sees them
  > The inner handler gets them: `WithAttrs` passes them on, which is all it does.
- [x] The pass and request IDs, since its records no longer pass through `passHandler`
  > The new logger is built on the handler `WithAttrs` returned, the inner one alone, so `passHandler`'s `Handle` never runs for it.
- [ ] Nothing, since `Handle` runs for every record whatever `WithAttrs` returns
  > `Handle` runs only on the handler a logger was built on, and this logger was built on the inner handler.

## Where does `pass_id` land in the record the last call writes?

- [ ] At the top level, after the `radio` group
  > Every attribute added after `WithGroup`, including those `Handle` adds, goes inside the group.
- [ ] Nowhere, since records in a group skip `Handle`
  > Every record goes through `Handle`; the group changes only where its attributes are written.
- [x] Inside `radio`, beside `margin_db`
  > An attribute added in `Handle` lands inside any group opened with `WithGroup`, which is why Apogee writes groups with `slog.Group`.
