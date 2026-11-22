As you walk through `seqs` in order, how can you tell whether you've already met the number in front of you?

Keep a `map[int]bool` of the numbers you've met and a new slice for the result. Append a number only when the map says `false` for it, and mark it `true` as you do.
