---
slug: a-table-of-commands
title: A Table of Commands
type: code_tests
position: 2
difficulty: 5
topics: functions, maps, sorting
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Nkem reads `dispatch` and asks one question: "Where does the twentieth command go?" Each new command should be one entry in a table, not another `if`.

`s.commands` maps a name to a `command`: its help text and the function that runs it. A comma-ok lookup also says whether the key was there:

```go
masks := map[string]float64{"SQT": 3}
m, ok := masks["CWR"]
fmt.Println(m, ok) // 0 false
```

## Your task

1. In `dispatch`, look the first word up in `s.commands`. Call a known command's `run` with `s`, the words after the name and `out`, and return its error; answer an unknown one as before.
2. Write `cmdHelp`, a function like `run`: a line per command, sorted by name, the name padded by `%-8s`, two spaces, then its help text. It ignores any words after `help`.
3. Add `help` to the table in `commands.go`, so that it answers:

```text
help      list every command
```
