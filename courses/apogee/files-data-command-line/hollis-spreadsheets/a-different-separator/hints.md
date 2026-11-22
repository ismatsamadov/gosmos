`Comma` and `Comment` are fields of the reader, of type `rune`. Set them with character literals in single quotes before reading anything.

Without `Comment`, the first comment line becomes a record of one cell, and every five-cell record after it is then the wrong width.

Return `ReadAll`'s two results as they are: quoting mistakes, such as a quote that never closes, come back as its error.
