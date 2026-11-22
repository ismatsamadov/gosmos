---
slug: nothing-new-nothing-sent
title: Nothing New, Nothing Sent
type: code_tests
position: 1
difficulty: 7
topics: bytes, http, json
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
Nkem does the sum: forty screens asking every fifteen seconds is 160 copies a minute, of a Board that changes a few times an hour.

An `ETag` header names one version of an answer. A screen that keeps the answer sends the name back in `If-None-Match`, and `304 Not Modified`, with no body, tells it that its copy is still current. The name must change only when the data does: `Refresh` numbers every view it builds, even one identical to the last, so every refresh looks like news.

## Your task

1. In `Refresh`, when the new view's JSON equals `b.json` (`bytes.Equal`), keep the view you have, with its number.
2. In `serveJSON`, send the view's version, `b.view.version`, as an ETag in quotes. When `If-None-Match` equals it, answer 304 with the ETag and no body:

```text
ETag: "v2"
```
