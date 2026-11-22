`strings.Fields(line)` gives you a slice of words. Check its length before you look at any word: a blank line gives an empty slice, and reading index 0 of an empty slice stops the program with a panic.

When there are no words, return `nil` before anything is written, so a blank line leaves only the next prompt. Otherwise the command is the word at index 0.

The unknown-command line keeps its format exactly. Only what goes between the quotes changes, from the whole line to its first word.
