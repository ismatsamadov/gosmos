Start with a position that cannot be a column, such as -1, and loop over the header with both the index and the cell.

When a cell says `name` and you have already found one, the header is ambiguous: return an error straight away. Otherwise remember the index.

After the loop, a position still at -1 means no cell matched. Put `name` in both errors, for example with `%s`, so Hollis knows which column to fix.
