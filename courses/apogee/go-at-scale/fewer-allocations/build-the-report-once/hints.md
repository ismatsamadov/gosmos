Keep the heading line as it is: it is built once. The `+=` inside the loop is what copies the report over and over.

Before writing anything, loop over the contacts once to add up the lines' lengths: each is `len(c.ID)` plus the length of the dish part and the newline. Then call `b.Grow` with the total and write the heading and each line with `WriteString`.

If the test still counts about ten allocations, the builder grew as it went: `Grow` must come before the first write, with the whole report's length. If it counts one or more a line, something in the loop allocates, as `fmt.Fprintf` does: write the parts with `WriteString`.
