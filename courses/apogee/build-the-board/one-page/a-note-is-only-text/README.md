---
slug: a-note-is-only-text
title: A Note Is Only Text
type: code_tests
position: 2
difficulty: 7
topics: security, templates
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Nkem tries your page on Thursday night's schedule, and the note on the Cape Wren pass comes out in bold. Hollis, Furrowlight's ops manager, wrote it as `<b>RUSH</b> & PRIORITY`, and the page took those tags as its own.

`html/template` escapes each value for the place it lands, so text a customer typed stays text. A `template.HTML` value is the exception: it means "trusted HTML, print it as it is", and nothing a customer typed is trusted.

```go
t := template.Must(template.New("").Parse(
	"<p>{{.}}</p>"))
t.Execute(os.Stdout, "wind > 72 & rising")
// <p>wind &gt; 72 &amp; rising</p>
```

## Your task

`board.html` passes each note through `asHTML`, a function in `page.go` that turns it into `template.HTML`. Take `asHTML` out of the template, the `FuncMap` and `page.go`, so a note prints as text:

```text
&lt;b&gt;RUSH&lt;/b&gt; &amp; PRIORITY
```
