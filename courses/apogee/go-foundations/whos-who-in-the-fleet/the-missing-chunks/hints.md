Ranging over `chunks` visits only the chunks that arrived. What could you loop over that visits every number the file should have, in order?

Range over the integer `total` itself, which counts from 0 to `total-1`. For each number, the comma-ok form tells you whether it is a key in `chunks`, even when its chunk is empty, and `append` collects the numbers that aren't.
