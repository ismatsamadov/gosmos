After `draft := live[:2]`, how many arrays hold the AOS times: one or two?

`slices.Clone` returns a new slice whose elements live in an array of their own. It is in the `slices` package.

Clone only the part the draft needs: `draft := slices.Clone(___)`.
