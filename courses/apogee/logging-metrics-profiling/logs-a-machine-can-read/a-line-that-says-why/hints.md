`logger` is already made at the top of `main`. Its `Error` method takes the message first, then the attribute's key and its value, each a string.

Put your call where `log.Print("error")` was. Then delete the `_ = logger` line, and the `"log"` import once nothing in `main.go` uses it.

If Submit finds no line at all, the old call is still the only one, and it writes to standard error. If your line differs, compare the message, the key and the value with the task letter by letter: `pass failed`, `site`, `LAN`.
