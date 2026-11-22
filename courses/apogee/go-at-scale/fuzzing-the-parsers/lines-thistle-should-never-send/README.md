---
slug: lines-thistle-should-never-send
title: Lines Thistle-2 Should Never Send
type: code_tests
position: 7
difficulty: 9
topics: errors, strings, testing
created: 2026-11-11T12:51:43Z
modified: 2026-11-12T08:43:36Z
---
Tuesday 17 April. Cirrine's twelve new satellites leave their rocket at 09:40 and send their first telemetry minutes later, and Mirela signs off on the parser only when every one of Joss's odd lines is refused and every good line comes back unchanged.

A strict parser checks each field's form, not only whether `strconv` can read it. Its last check is your round trip: a line that isn't `formatLine` of its own frame is refused, so whatever passes is written the one way tlmd accepts. Each error names the field and the form it wants:

```go
code := "pel"
if strings.ToUpper(code) != code {
	return fmt.Errorf("%q: want capitals", code)
}
```

## Your task

Make `parseLine` strict:

1. Refuse a satellite ID that isn't 1 to 24 lowercase letters, digits and `-`, starting with a letter, and a mode other than the four.
2. Read the time with the layout `2006-01-02T15:04:05Z`, and readings written with digits, a point and a sign only.
3. Last, refuse a line `formatLine` doesn't write back exactly: `line: not canonical`.

> [!NOTE]
> The four messages `TestParseLine` checks are in its table, word for word.
