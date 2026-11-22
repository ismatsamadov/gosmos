`switch format` with one `case` for `"json"` and one for `"csv"` covers both files. Each case can return the reader's two results directly.

The `default` case is every other format, the empty string included: return nil and an error built with `fmt.Errorf` and `%q`, so the message shows the name Paz typed.

Never look at the input to decide: the test hands a CSV sheet to `"json"` and expects the JSON decoder's error, not a quiet switch to the other reader.
