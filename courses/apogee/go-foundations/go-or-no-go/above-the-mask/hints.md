Look at the tests: a pass that peaks exactly at its mask still clears it, and not every station's mask is 5°. Which comparison includes the equal case, and what does `maxElevDeg` need comparing with?

`>=` means "at or above". Inside each block, assign with `=`, since `call` already exists: `:=` would declare a second `call` inside the block that is never used, and Go won't build that.

Shape: `if maxElevDeg ___ maskDeg {` sets `call` to `___`, then `} else {` sets it to `___`, and a closing `}` comes before the `return`.
