Look at the minutes in order, 0 to 14: the first one that reaches the mask is the answer. What should the function return if the loop finishes without finding one?

Keep the answer in a variable that starts at -1. Inside a `range` loop, an `if` with `>=` spots the first minute that clears the mask; save that minute, then `break`.

Shape: `aos := -1`, then `for minute := range ___ {` holding an `if` that sets `aos` and breaks, then `return aos` after the loop.
