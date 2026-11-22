`cmp.Or` takes the comparisons in the order they count and returns the first one that isn't zero: the AOS comparison first, then the pass IDs'.

`strings.Compare(a.ID, b.ID)` compares two pass IDs the way `a.AOS.Compare(b.AOS)` compares two times, with −1, 0 or +1. Give it to `cmp.Or` after the AOS, and drop the line with the TODO.

If a test says Pipit-3 came before Brine-3, the IDs aren't compared yet. If it says an earlier pass came second, the ID is the first key: swap the two.
