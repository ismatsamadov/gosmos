Declare the pool as a package-level `var` with a `New` field: a function that makes the 512-pass slice, as `getScratch` does now, and returns a pointer to it.

`Get` returns `any`, so assert it to `*[]PassView`. A scratch may come back holding the passes of its last station: set its length to 0 through the pointer before handing it out.

In `putScratch`, compare `cap(*p)` with `maxScratch` and `Put` only a scratch that isn't bigger. If the tests never get the 65,536-pass scratch back, check that comparison's `=`.
