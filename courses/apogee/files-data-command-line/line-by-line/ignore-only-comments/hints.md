Make a second view of the line with `strings.TrimLeft` and the characters space and tab, and ask your two questions of that view instead of the line.

An empty view means a blank line; a view starting with `#` means a note. Append the original line, not the trimmed view: the tests expect indented readings back unchanged.
