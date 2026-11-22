---
slug: wind-on-every-panel
title: Wind on Every Panel
type: code_tests
position: 4
difficulty: 7
topics: conditionals, loops, templates
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Ingrid wants each station's wind beside the passes, and the name of whoever is on call. "At 04:00, nobody should have to ask."

`{{define "name"}}` names a part of a template, and `{{template "name" .}}` runs that part with the value you give it, here once for each shift `{{range}}` visits:

```html
{{define "shift"}}<li>{{.Name}}</li>{{end}}
<ul>{{range .Rota}}{{template "shift" .}}
{{end}}</ul>
```

`{{if}}` and `{{else}}` choose between two parts, and a nil pointer counts as false.

## Your task

In `board.html`, after the passes:

1. Define `panel`, which writes one `<li>` on one line. Its class is `panel` and the panel's `Call`, as in `class="panel warn"`. Its text gives the wind and gusts with no decimals, or says there is no reading:

```text
YST wind 41 km/h, gusts 55
SQT no reading
```

2. In a `<ul class="weather">`, write the panel of each station in the view's `Weather`, one per line.
3. Add the on-call line:

```text
<p class="oncall">On call: Nkem Halloran</p>
```
