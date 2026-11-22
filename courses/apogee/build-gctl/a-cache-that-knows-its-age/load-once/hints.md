`passes` turns on `c.loaded`. While it is false, call `c.load()`; on success set both `c.cached` and `c.loaded`, and on an error return `nil` and that error without setting either.

In `newShell`, the new field is a pointer to a `passCache` whose `load` is `loadTonight` itself, written without parentheses: the cache calls it when it needs to, not `newShell`.

In `cmdPasses`, ask the cache after the argument checks and before cloning and sorting. The cache hands you its own slice, so keep cloning before you sort.
