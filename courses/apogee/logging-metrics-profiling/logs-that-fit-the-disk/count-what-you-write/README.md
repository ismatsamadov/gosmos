---
slug: count-what-you-write
title: Count What You Write
type: code_tests
position: 1
difficulty: 7
topics: files, interfaces, streams
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Nkem's first question about the new writer: "How big is the file now, without asking the disk?" To start a new file at 64 MiB, it must know after every line, so it counts what it writes, from the file's size when it opened it.

`Write` returns `n`, how many bytes of `p` it took, and an error. When it fails partway, as on a full disk, `n` is how many bytes did reach the file, so a count of `n`, not of `len(p)`, stays equal to the file's size. A writer can wrap another and add to what it does:

```go
type tally struct {
	w     io.Writer
	calls int
}

func (t *tally) Write(p []byte) (int, error) {
	t.calls++
	return t.w.Write(p)
}
```

## Your task

`countingWriter`'s `Write` passes each write on and counts nothing. Make it add the `n` the inner writer returns to `c.n`, whether that write failed or not, and return `n` and the error as they came.
