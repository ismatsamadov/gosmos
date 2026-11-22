Start at the end: `if formatLine(fr) != line`, return `fmt.Errorf("line: not canonical")`, where `fr` is the frame just built. That alone refuses most of Joss's lines.

The round trip can't see a capital in a satellite ID or an unknown mode, since `formatLine` writes both back as they came: check `f[0]` rune by rune, and the mode against the four. Parse the time with the layout `"2006-01-02T15:04:05Z"`, so `+00:00` gets the time's own error.

If a message in `TestParseLine` differs from yours, copy the `want` beside it: each puts the field's text in `%q`, such as the sequence's `"-1"`.
