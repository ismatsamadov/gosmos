---
slug: saying-goodbye
title: Saying Goodbye
type: code_tests
position: 3
difficulty: 5
topics: errors, loops
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Ingrid types `exit`, and gctl calls it an unknown command. A session must end when she says so, and if her terminal drops halfway, gctl must report it instead of saying goodbye.

A sentinel error is one error value, declared once, that callers recognize with `errors.Is`. So far `Run` has dropped what `dispatch` returns. Now it tells three endings apart: `errExit` ends the session, another command's error is printed while the session goes on, and a failed read is returned.

## Your task

1. In `commands.go`, declare `errExit` with `errors.New`, write `cmdExit` returning it, and add `exit` to the table.
2. In `Run`, leave the loop on `errExit`, and print other errors as `error: <message>`.
3. After the loop, return the scanner's `Err` if it has one; otherwise write `bye` and a newline.

`help` then shows:

```text
exit      leave the shell
help      list every command
```
