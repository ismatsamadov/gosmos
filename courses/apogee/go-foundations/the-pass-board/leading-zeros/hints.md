What should `clock(4, 7)` return, and which part of a verb pads a number with zeros? The lesson's example pads to five digits.

Use `fmt.Sprintf`, which returns the string instead of printing it, with a zero-padded verb of width 2 for each number.

The format is `%__:%__`: two verbs with a colon between them and no `\n`, followed by `h` and `m`.
