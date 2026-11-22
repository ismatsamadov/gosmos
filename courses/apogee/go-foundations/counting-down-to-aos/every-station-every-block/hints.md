The station loop already prints each row's code and ends the row. What has to happen 12 times between those two lines?

A second `for`, over the 12 blocks, goes between the code and `fmt.Println()`. Inside it, `if` and `else` choose the character, printed with `fmt.Print` so it stays on the row.

Shape: `for block := range ___ {`, then `if visible(___, ___) {` printing `#`, `} else {` printing `.`, `}`, and a final `}` to close the loop.
