---
slug: the-fuzzers-find
title: The Fuzzer's Find
type: code_tests
position: 0
difficulty: 8
topics: errors, strings, testing
created: 2026-11-11T12:51:44Z
modified: 2026-11-12T08:43:35Z
---
The line Thistle-2 sent at 05:10 on Saturday had four fields where tlmd expects six. Mirela, Apogee's security engineer: "A parser may refuse a line. It may not crash on one."

Overnight, she fuzzed the parser. A fuzz test, such as `FuzzParseLine(f *testing.F)`, runs one check on many inputs. Plain `go test` gives it its seed corpus: the values `f.Add` adds, and every file under `testdata/fuzz/FuzzParseLine/`, where she saved Saturday's line. A corpus file holds one input, as this one for another fuzz test does:

```text
go test fuzz v1
string("pel")
```

## Your task

`parseLine` indexes its fields without counting them, so Saturday's four-field line panics at `f[4]`. Before it reads any field, refuse a line that hasn't exactly six, with the line's own count in this error:

```text
line: want 6 fields, got 4
```

> [!TIP]
> Run `go test -fuzz=FuzzParseLine -fuzztime=30s` in a terminal to invent lines for 30 seconds. Any that fails is saved in that folder, where plain `go test` runs it from then on.
