Count the characters in the target table: each name and its padding fill 13 columns, then comes one space, then 5 columns for the mask.

Names follow the pattern of the example's `%-6s` and masks that of its `%6.2f`, each with your own width and decimals. `MASK` in the header is text, so its verb is `s`, as wide as the masks below it.

Each row's format has the shape `"%___s %___f\n"` and the header's `"%___s %___s\n"`, with the same widths in both.
