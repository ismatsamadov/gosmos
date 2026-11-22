Only the line that declares `out` changes: the loop and the `append` stay as they are.

`make` for a slice takes a length and a capacity. The slice should start empty, with room for `n` passes.

If the test counts more passes than it handed in, the slice started with `n` zero values: its length must be 0, and `n` goes in its capacity.
