Inside the loop, before computing a step's code, `continue` when `step <= c.last[user]`: that step, or a later one, has already been used.

Record the step at the moment a code matches, `c.last[user] = step`, right before returning true. Recording `now` instead would let a fast phone's code, the step after, count twice.

If Mirela's replay still passes, the step check is missing; if a wrong code spoils the right one, `c.last` changes before the comparison succeeds.
