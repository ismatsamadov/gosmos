In `react`, which dish does `d.Stow()` stow: the one in `dishes`, or a copy that's gone after each turn of the loop?

Loop over the indexes of `dishes`, as the example's second loop does with `stations`, and call `Stow` and `Unstow` on the element at each index instead of on a copy.
