The loop calls `checkBattery(v)` and drops what it returns. What should you keep, and what should you check before printing `stored`?

After `if err != nil`, print the rejection with `%v` for the error, then use `continue` so both the stored line and the count are skipped.
