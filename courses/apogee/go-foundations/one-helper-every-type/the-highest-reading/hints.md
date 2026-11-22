Where should the search start? If it starts from the zero value, what does `Max` return when every reading is below zero?

Check `len(items)` first. Then start from `items[0]` and compare each item with `>`, which `cmp.Ordered` allows.

The shape: return the zero value and `false` for an empty slice; otherwise keep a `best` that a `range` loop replaces whenever an item is greater, then return `best` and `true`.
