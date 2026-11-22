Compare `len(old)` and `len(new)` first, and return the first form with `fmt.Errorf` when they differ. Add `fmt` to the imports.

Then loop over the indexes of `old`, and at the first index where the two contacts' `ID` or `Dish` differ, return the second form with that index, then the old contact's ID and dish, then the new one's.

Only the first difference counts: return as soon as you find it. Indexes count from 0, as Go's do, and a contact whose ID matches but whose dish doesn't is a difference too.
