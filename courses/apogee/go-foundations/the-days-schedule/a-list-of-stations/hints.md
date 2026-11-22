If the first code is at index 0, at which index is the fifth? How could you get that number from the length?

A slice literal is its type, then the values in braces: `[]string{...}`. `len(codes)` counts the elements.

The last index is always one less than the length, however many stations there are.
