`range` counts from 0, but the log counts from 1. What do you do to the index before printing it?

A `range` loop over `ids` with two loop variables gives you each index and pass ID in turn. The print goes inside the loop's braces.

`fmt.Println` puts one space between its arguments, which is exactly the format the log wants.
