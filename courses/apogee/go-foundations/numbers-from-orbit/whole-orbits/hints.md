How many times does a 93-minute orbit fit into 1,440 minutes, and how many minutes are left after the last full one? That is a division and its remainder.

Both variables are `int`s, so `/` gives the whole number of orbits and `%` gives the minutes the division leaves over. Divide the day by the period, not the other way round.

Skeleton: `fullOrbits := minutesPerDay ___ periodMin` and `minutesLeft := minutesPerDay ___ periodMin`.
