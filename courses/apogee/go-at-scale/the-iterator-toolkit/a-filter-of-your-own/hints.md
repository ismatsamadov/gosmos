Return a function that takes `yield`, as `Repeat` does. Inside it, range over `seq` and ask `keep` about each value.

Call `yield` only for a value `keep` returns true for. When that call returns false, return at once: the loop over `seq` ends, and `seq` stops too.

If a test says the schedule yielded 260 contacts for a loop that stopped at Tirrin Rock's second, your loop went on after `yield` returned false, or read the whole schedule before yielding; a panic about "continued iteration" means it called `yield` again.
