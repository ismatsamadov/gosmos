What does `battery` do with the second result of `strconv.ParseFloat` right now?

Keep that result in `err` and check it with `if err != nil` right after the call, the way you checked every error in the last chapter.

`fmt.Errorf` with `%w` puts the parse error inside yours: write the field's name, `batt_v`, before the `%w` in the message.
