---
slug: words-not-lines
title: Words, Not Lines
type: code_tests
position: 1
difficulty: 4
topics: slices, strings
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Ingrid types `track`, two spaces, then `pipit-1`, and gctl calls the whole line an unknown command. "Negative. At 04:00 nobody counts spaces."

`strings.Fields` splits text into words at any run of spaces or tabs, and never returns an empty word. A blank line gives no words at all:

```go
words := strings.Fields(" stow\tCWR  wind")
fmt.Println(len(words), words[0]) // 3 stow
fmt.Println(len(strings.Fields(" \t "))) // 0
```

## Your task

In `dispatch`, split the line with `strings.Fields`; the first word is the command. A line with no words writes nothing, and an unknown command names only its first word, so Ingrid's line gets:

```text
unknown command "track"; type help
```

> [!TIP]
> Add `"strings"` to the imports of `shell.go`.
