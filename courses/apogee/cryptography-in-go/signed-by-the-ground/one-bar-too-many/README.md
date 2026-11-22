---
slug: one-bar-too-many
title: One Bar Too Many
type: code_tests
position: 1
difficulty: 7
topics: crypto, errors, strings
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Tuesday. On the whiteboard Mirela writes two commands: `SET_MODE` with the one argument `SAFE|REBOOT`, and `SET_MODE` with two, `SAFE` and `REBOOT`. Both become the same signed text, `pipit-2|118|SET_MODE|SAFE|REBOOT`, so one signature fits both.

The exact text that gets signed, the canonical form, must have one meaning only, so any field that could be read two ways is refused. Checking a field is a loop over its characters:

```go
for i, r := range "YsT" {
	if r < 'A' || r > 'Z' {
		fmt.Println(i, string(r)) // 1 s
	}
}
```

## Your task

`canonical` joins the fields with `|` and checks nothing. Make it refuse any field these rules don't allow, an empty one included, with an error that starts with the field's name, as in `argument 1 contains "|"`:

```text
satellite ID     lowercase letters, digits, -
sequence number  at least 1
operation        capital letters and _
argument 1, 2…   1 or more of ! to ~, never |
```
