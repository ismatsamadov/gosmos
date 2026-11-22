---
slug: values-behind-the-fields
title: Values Behind the Fields
type: code_tests
position: 3
difficulty: 9
topics: formatting, pointers, reflection, structs
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T03:14:43Z
---
Wednesday 30 May. When passplan 2 starts, it logs its settings, and the hand-written log line names six of `Config`'s eight fields: nobody added `RunLimit` and `Quiet`.

`reflect.ValueOf(v)` holds the value itself, and a struct's `Value` has `Fields()`, which ranges over pairs: each field's `StructField` and its `Value`, which `fmt` prints as the value it holds. A pointer's `Value` has the kind `reflect.Pointer`, and its `Elem()` is what it points to. Reflection reads an unexported field but won't hand its value out: `Interface()` on it panics, so check `StructField.IsExported` first.

```go
w := Window{From: dawn, To: dusk}
v := reflect.ValueOf(&w)
if v.Kind() == reflect.Pointer {
	v = v.Elem()
}
to := v.FieldByName("To")
fmt.Println(to.Type()) // time.Time
```

## Your task

`describe` returns `nil`. Return `<name>=<value>` for each exported field of `v`'s struct, in order, with the value printed by `%v`. When `v` is a pointer, describe the struct it points to; anything that isn't a struct gives `nil`.

> [!NOTE]
> The tests include a dish loader's settings, whose password for its host's line is unexported and must stay out of the log.
