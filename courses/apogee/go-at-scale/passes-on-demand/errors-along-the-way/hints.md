Where the loop now says `continue`, call `yield` with `PassView{}` and the error, then return: an archive with a damaged line ends there.

Wrap the error with `fmt.Errorf`: the line's number in the message, and `%w` where the error goes. `n` already counts lines from 1. Remember to import `fmt`.

`sc.Scan` returns false at the end of the input and when a read fails; after the loop, `sc.Err()` tells the two apart. The test reads a damaged third line, a damaged fifth line, and a disk that fails after two lines.
