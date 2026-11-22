---
slug: page-one-is-wrong
title: Page One Is Wrong
type: code_tests
position: 0
difficulty: 6
topics: errors, networking, streams, strings
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Monday 13 March. dishd, the dish service, now reads a controller's first line as page one of the manual describes it, and refuses a controller that is fine. Sunniva: "Simulator fine. dishd says no. Page one wrong."

`br.ReadString('\n')` reads through the first `\n` and returns the line with its ending. Remove exactly the bytes that end it and nothing more: a stray space or a second `\r` means the line isn't what the protocol says.

```go
br := bufio.NewReader(strings.NewReader(
	"Wind 64\r\nGusts 81\r\n"))
line, err := br.ReadString('\n')
fmt.Printf("%q %v\n", line, err)
// "Wind 64\r\n" <nil>
```

## Your task

`readBanner` removes only the `\n`, as page one says, so the real banner's `\r` stays on the end and the version reads `"1\r"`. Replace the `TODO` so a `\r` right before the `\n` goes too, and nothing else.

> [!WARNING]
> `strings.TrimSpace` would also accept `PROTO 1 ` with a space at the end, a banner no controller sends.
