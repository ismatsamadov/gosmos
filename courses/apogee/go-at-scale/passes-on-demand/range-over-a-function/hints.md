Inside the function `Days` returns, loop once for each of the `n` days; `for i := range n` counts `i` from 0 to `n-1`.

Day `i` is `from` moved on by `i` whole days, which `AddDate` does. Hand each day to `yield`, and return if it gives back false.

If the week starts on Mon 2 Jul, the first day is one too late: day 0 is `from` itself. If nothing prints at all, the loop never calls `yield`.
