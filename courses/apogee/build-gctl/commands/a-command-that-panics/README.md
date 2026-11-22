---
slug: a-command-that-panics
title: A Command That Panics
type: code_tests
position: 4
difficulty: 6
topics: defer, errors, functions
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Nkem asks one question: "A command has a bug and panics at 04:12. What does Ingrid see?" Today gctl ends in the middle of the pass, and the night shift loses its prompt.

A deferred function runs as its function returns, even through a panic, and `recover` called there stops the panic. With a named result, the deferred function can also change what is returned:

```go
func div(a, b int) (q int, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("cannot divide")
		}
	}()
	return a / b, nil
}
```

## Your task

In `dispatch`, name the error result, and recover a command's panic so that `dispatch` returns an error naming the command. `Run` then prints it, and the session goes on:

```text
error: internal error in doppler
```
