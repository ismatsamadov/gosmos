Cut the line that calls `stationDay` and put it above the loop. Then `cwr` and `day` are made once, and the loop's body uses them on every turn.

Leave the call to `planDay` inside the loop: it's the work being measured, and the tests still want it 20 times, on all of Cape Wren's 360 candidates.

If `AllocsPerOp` is still in the thousands, something in the loop still predicts: only the line that calls `planDay` belongs there.
