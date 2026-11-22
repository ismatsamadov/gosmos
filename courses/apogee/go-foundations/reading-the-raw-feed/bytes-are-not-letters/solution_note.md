`len` counts bytes, and `ñ` and `á` take two bytes each in UTF-8; `utf8.RuneCountInString` counts runes, so each name reports the number of characters a reader sees.
