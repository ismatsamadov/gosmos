Convert the receiver back to `time.Duration` to use its `String` method, then hand that text to `json.Marshal`, which adds the quotes.

In `UnmarshalJSON`, decode `data` into a Go `string` first. That refuses numbers and every other JSON type, and turns escapes back into characters; then parse the string.

Give `UnmarshalJSON` a pointer receiver and assign the parsed value through it at the end. A method on a plain `Delay` would set a copy that is thrown away.
