---
slug: help-is-a-success
title: Help Is a Success
type: code_tests
position: 2
difficulty: 6
topics: errors, programs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Ingrid's night shift typed `gctl passes -h` to see its options. The usage came out on standard error with status 2, and their checklist script logged a failure.

On `-h` or `-help`, `Parse` writes the usage to the flag set's output and returns `flag.ErrHelp`. Help that was asked for is a success: usage on standard output, status 0. Any other parse error is a mistake: message and usage on standard error, status 2. `Parse` writes before you know which, so give it a buffer, then send the buffer where it belongs.

```go
err := fs.Parse([]string{"-help"})
if errors.Is(err, flag.ErrHelp) {
	fmt.Println("help was asked for")
}
```

## Your task

Implement `parseFlags`: point `fs`'s output at a `bytes.Buffer`, then parse `args`. After help, copy the buffer to `stdout` and return 0 and false; after any other error, copy it to `stderr` and return 2 and false. Otherwise write nothing and return 0 and true: the command goes on.
