---
slug: name-the-bad-line
title: Name the Bad Line
type: code_tests
position: 6
difficulty: 4
topics: errors, streams, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T17:07:42Z
---
Paz runs the elevation reader over the whole log, and it stops with a syntax error from `ParseFloat` for the text `41..70`. The typed reading is somewhere in 200,000 lines, and the error does not say where.

Count lines as you scan, from 1, and count every line, the night shift's `#` notes and blank lines included, so the number matches the line an editor shows. Then wrap the error with that number using `%w`: the message gains its place, and `errors.Is` still finds the original error inside.

```go
eof := io.ErrUnexpectedEOF
err := fmt.Errorf("frame %d: %w", 1207, eof)
fmt.Println(err) // frame 1207: unexpected EOF
ok := errors.Is(err, eof)
fmt.Println(ok) // true
```

## Your task

In `readElevations`, count each line before anything can skip it, and wrap a parsing error with `%w` and the line's number. A missing column on line 1 must read:

```
line 1: 3 columns instead of 4
```

The read-only `parse.go` holds the column and elevation helpers from the last two lessons.
