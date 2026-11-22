---
slug: quick-check-templates
title: Quick Check: Templates
type: quiz
position: 3
difficulty: 6
topics: security, templates
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Before the page goes up on a second screen, the kitchen's, Nkem has three questions about lines like these:

```html
<li class="panel {{.Call}}">{{.Station}}</li>
<script>var call = {{.Call}};</script>
```

```go
var buf bytes.Buffer
err := b.page.Execute(&buf, data)
```

## Your task

Answer three questions: how one value is escaped in two places, what `template.HTML` changes for a note, and why the page is executed into a buffer first.
