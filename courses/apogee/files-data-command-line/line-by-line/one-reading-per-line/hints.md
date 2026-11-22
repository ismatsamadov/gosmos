`Scan` makes a good loop condition: it returns true while another line is waiting, so a loop on it runs once per line.

Add one to `n` in the loop's body without looking at the text. An empty line counts as much as a reading, and a last line without a newline still makes `Scan` return true once.
