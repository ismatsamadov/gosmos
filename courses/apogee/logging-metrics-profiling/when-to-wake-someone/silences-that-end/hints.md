Loop over the silences and return true at the first one that holds the station back; after the loop, return false. A silence for another station never counts.

`at` is inside when it is not before `From`, so the start itself counts, and when it is before `Until`, so the end doesn't: `Before` gives you both checks.

If 11:00 is still held, the end is included: a silence until 11:00 is over at 11:00. If 09:00 isn't held, the start is left out.
