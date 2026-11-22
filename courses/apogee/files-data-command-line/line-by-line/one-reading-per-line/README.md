---
slug: one-reading-per-line
title: One Reading per Line
type: code_tests
position: 0
difficulty: 4
topics: loops, streams, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
Paz brought down the file where the computer steering Salar Quieto's dishes writes one line per pointing reading: the controller log. The station counted 200,000 lines, and Paz wants the van's copy to give the same count.

`bufio.Scanner` reads a stream one line at a time through a small buffer, so the whole file never sits in memory. Each `Scan` moves to the next line and returns false when none is left; `Text` returns the line without its newline:

```go
r := strings.NewReader("PEL\nSQT\n")
sc := bufio.NewScanner(r)
for sc.Scan() {
    fmt.Println("station", sc.Text())
}
```

## Your task

Replace the TODO in `countLines` with a loop that counts every line, blank ones and a last line without a newline included. The supplied `return` already hands back `sc.Err()`; the next lesson shows why.
