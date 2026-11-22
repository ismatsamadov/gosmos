`fmt.Fprintln` writes to whatever `io.Writer` it receives, so the same `writeBoard` draws on the terminal through `os.Stdout` and into a `strings.Builder` in the test.
