A `switch` on `sev` with three cases fits: `Warning` and `Info` each set `a.To` to a fixed text, and `Critical` looks through the rota.

The rota is oldest first: loop over it and set `a.To` to the `Name` of each shift whose `From` is not after `a.At`, so the last one you set is the latest. Return `a` at the end.

If Monday 08:00 still goes to Tomás, the comparison leaves out a shift that starts exactly at the alert's time: a shift is on from its first minute.
