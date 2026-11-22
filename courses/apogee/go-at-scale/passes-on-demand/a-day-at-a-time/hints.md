Put the prediction inside a loop over the days of `w`, as in `Passes`: each day runs from `from` to a day later, or to `w.To` if that comes first.

For each day, collect the passes of every satellite at every station that can serve it, predicting that day only, and sort that day's slice with `byAOS` before yielding any of it.

Return once `yield` returns false. The test breaks after 10 passes and fails if `predict` had been asked about anything but 1 July by then.
