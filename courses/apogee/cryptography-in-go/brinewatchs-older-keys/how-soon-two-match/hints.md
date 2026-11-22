A `map[string]int` from a code to the number of the text that had it remembers every code so far. Make it before the loop, and start the loop at text 0, so the first text's code goes in too.

For each text, look its code up with the two-value form, `j, found := seen[code]`. If it is there, return `j`, this text's number and true; if not, store this text's number under the code and go on.

The 2-byte line should show a match within a few hundred tries and the 3-byte line within a few thousand. If the second line still says no two match, the loop still compares every text with text 0 alone.
