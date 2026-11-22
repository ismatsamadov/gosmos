---
slug: a-page-from-a-template
title: A Page From a Template
type: code_tests
position: 0
difficulty: 6
topics: bytes, http, templates
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
The Board, the screen in the Fishbowl ops room that lists the next passes, is a laptop plugged into a television. When it rebooted during last week's night shift, the Fishbowl went twenty minutes without a Board. Ingrid wants a program to serve it.

`page.go` builds the page's template into the program, and `NewBoard` parses it once into `b.page`, an `html/template` template:

```go
//go:embed board.html script.html
var pages embed.FS
```

A template writes as it executes, so one that fails halfway has already sent half a page, with status 200. Execute it into a `bytes.Buffer` instead, and nothing leaves until the page is whole:

```go
var buf bytes.Buffer
err := t.Execute(&buf, data)
```

## Your task

`servePage` answers `GET /` with `TODO`. Make it:

1. Execute `b.page` into a buffer, with a `pageData` of `b.current()` and `b.now()`.
2. If that fails, answer 500 with `internal error` and none of the page.
3. Otherwise, set `Content-Type` to `text/html; charset=utf-8` and send the buffer. At 03:10 the page's heading reads:

```text
<h1>The Board · 03:10 UTC</h1>
```

> [!NOTE]
> `script.html` is the page's script: given, and never checked.
