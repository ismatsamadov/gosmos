---
slug: a-prompt-that-waits
title: A Prompt That Waits
type: code_tests
position: 0
difficulty: 4
topics: loops, streams
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Ingrid wants the night shift, the team that runs the passes after midnight, to keep gctl open through a whole pass. First, gctl has to wait for her.

A `bufio.Scanner` reads one line per call to `Scan`, which returns `false` once the input has ended; `Text` gives the line without its newline:

```go
r := strings.NewReader("YST\nPEL\n")
sc := bufio.NewScanner(r)
for sc.Scan() {
	fmt.Println("station", sc.Text())
}
```

## Your task

`Run` in `shell.go` returns at once. Make it write the prompt `gctl> ` before each line, with no newline, pass each line to `s.dispatch`, and write the prompt once more at the end. For the lines `track` and `point`, `Run` then writes:

```text
gctl> unknown command "track"; type help
gctl> unknown command "point"; type help
gctl> 
```

> [!NOTE]
> `dispatch` returns an error that stays `nil` until Saying Goodbye, three lessons on.
