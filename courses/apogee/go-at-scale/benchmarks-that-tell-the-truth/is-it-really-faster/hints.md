Keep `n`, `cfg` and `friday` at the top, then call `b.Run("old", ...)` and `b.Run("new", ...)`. A slice of name and planner pairs, ranged, saves writing the inner part twice.

Inside each, range over 6, 12 and 24: build the window that ends that many hours after `friday` before the inner `b.Run`, name it `hours=6` and so on with `fmt.Sprintf`, and in its function call `b.ReportAllocs()` before the loop that plans. Add `fmt` to the imports.

If a count of allocations is off by a little, something besides the planner runs inside the loop; if by a lot, that sub-benchmark plans another window, another planner, or with other settings.
