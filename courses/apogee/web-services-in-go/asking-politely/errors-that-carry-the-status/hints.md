Inside the status check, read the start of the body first: `io.LimitReader(resp.Body, 512)` gives a reader that stops after 512 bytes, and `io.ReadAll` reads what it gives.

Build a `&StatusError{...}` from `resp.StatusCode`, `resp.Status` and the text, passed through `strings.TrimSpace` to drop the newline at its end. Then return it inside `fmt.Errorf` with `%w` where `%s` and the status were.

If `errors.AsType` finds nothing, the error is still formatted with `%s` or `%v`, which keeps only its text: only `%w` keeps the `*StatusError` itself inside the message.
