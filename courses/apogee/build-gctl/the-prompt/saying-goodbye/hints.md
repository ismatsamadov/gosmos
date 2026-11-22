Declare `errExit` once, outside any function in `commands.go`, so every `return errExit` hands back that same value and `errors.Is(err, errExit)` recognizes it in `Run`. Add `"errors"` to the imports of both files.

Keep what `dispatch` returns in a variable. On `errExit`, `break` out of the loop before writing another prompt; on any other error, write the `error:` line with `fmt.Fprintf` and `%v`, then carry on to the next prompt as usual.

After the loop, `sc.Err()` is `nil` both after `exit` and at the end of the input, and it holds the reader's error when reading failed. Check it before writing `bye`, so a broken input returns its error and never says goodbye.
