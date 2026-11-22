`range` over a string yields whole runes, so the two bytes of `ñ` arrive as one rune and become one `_`, and comparing each rune with ASCII ranges keeps only the characters a radio call sign allows.
