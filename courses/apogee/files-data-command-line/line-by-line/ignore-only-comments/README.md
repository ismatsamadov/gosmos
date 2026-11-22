---
slug: ignore-only-comments
title: Ignore Only Comments
type: code_tests
position: 7
difficulty: 4
topics: streams, strings, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
With line numbers, Paz finds the reader's next stop in a minute: line 198,114 has 6 columns instead of 4. It is a night-shift note indented with a tab, which the check for `#` at the very start of a line missed.

Classify a line by its first character that is not a space or tab: none means a blank line, `#` a note, anything else a reading. Trimming spaces and tabs from the left with `strings.TrimLeft` gives that view of the line. Use it to decide, and store the line as written.

```go
raw := "\t# LG warm after noon, Paz"
view := strings.TrimLeft(raw, " \t")
fmt.Println(strings.HasPrefix(raw, "#"))  // false
fmt.Println(strings.HasPrefix(view, "#")) // true
```

## Your task

Replace the check under the TODO in `dataLines`: skip a line of only spaces and tabs, or one whose first other character is `#`. Return every other line exactly as written, leading spaces and tabs included.
