Look at what `main` prints now, and at the constant above it: the number you need is already there.

`fmt.Printf` puts each value you pass it where a `%d` stands in its format string, in order.

`Printf` adds nothing on its own: end the format string with `\n` so the line ends as the
expected output does.
