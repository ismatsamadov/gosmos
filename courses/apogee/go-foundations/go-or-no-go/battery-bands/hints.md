A battery at 8.3 V is above all three limits, so several conditions are true for it. Since only the first true one counts, which limit should Go test first?

Chain the checks with `else if`, from the highest limit down, and let a final `else` set the band below the lowest limit. A limit itself belongs to the band it starts, so compare with `>=`.

Shape: `if battV >= ___ {` sets `band` to FULL, two `} else if battV >= ___ {` blocks set OK and LOW, and `} else {` sets the last band.
