Point the flag set's output at a `bytes.Buffer` with `SetOutput` before calling `Parse`; `Parse` fills it with any message and usage.

Test the error from `Parse` with `errors.Is` against `flag.ErrHelp` first. Then any other error is a mistake. The buffer's `WriteTo` copies what it holds to a writer.

When help reaches stderr with status 2, the help case is being handled as one more error: check for `flag.ErrHelp` before the general `err != nil` case.
