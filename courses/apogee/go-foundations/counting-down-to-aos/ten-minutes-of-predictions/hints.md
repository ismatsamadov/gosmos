How many lines does the expected output have, and which minute does it start at?

A `range` over the number 10 runs ten rounds and counts from 0, which gives exactly the minutes you need.

Shape: `for minute := range ___ {`, the `Printf` line inside, then `}`. Delete `minute := 0`, since the loop declares `minute` itself.
