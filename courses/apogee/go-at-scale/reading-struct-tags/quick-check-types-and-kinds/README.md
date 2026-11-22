---
slug: quick-check-types-and-kinds
title: Quick Check: Types and Kinds
type: quiz
position: 2
difficulty: 8
topics: pointers, reflection, types
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T03:14:43Z
---
Tuesday 29 May. On the whiteboard beside the backend desks, Nkem writes a question: is `Turnaround` a `time.Duration` or an `int64`? Under it, five lines read passplan 2's settings through a pointer.

```go
cfg := Config{Turnaround: 2 * time.Minute}
t := reflect.TypeOf(&cfg).Elem()
f, _ := t.FieldByName("Turnaround")
fmt.Println(f.Type, f.Type.Kind())
fmt.Println(f.Tag.Get("unit") == "")
```

The fourth line prints `time.Duration int64`, and the fifth prints `true`: `Turnaround` has no `unit` key.

## Your task

Answer Nkem's three questions: why a rule needs the type as well as the kind, what `Get` says about a missing key, and what the `Elem()` is for.
