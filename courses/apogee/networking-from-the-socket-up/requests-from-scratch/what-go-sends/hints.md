Put the `ReadString`, its error check and the `Printf` inside a `for` loop with no condition, so each turn reads and prints one line.

The request's head ends with a line that holds nothing but CR LF. Compare each line with `"\r\n"` after printing it, and leave the loop with `break` when it matches.

If the program prints five lines and then an i/o timeout, the loop kept reading after the blank line: Go's client sends nothing more until it has an answer, so stop at the blank line, and the given code answers.
