---
slug: defaults-from-tags
title: Defaults from Tags
type: code_tests
position: 7
difficulty: 9
topics: pointers, reflection, structs
created: 2026-11-12T03:14:42Z
modified: 2026-11-12T09:50:21Z
---
Tuesday 5 June. Friday's `Slots: 0` was a setting left out, not a wish, so Nkem decides that zero means "use the default", which every field of `Config` already gives in a `default` tag.

To set a struct's fields, reflection needs the struct itself, not a copy: `reflect.ValueOf(s)` holds a copy, while the `Elem()` of a pointer's `Value` is the struct itself. `CanSet` says whether a `Value` can be set: true for that struct's exported fields, false for a copy or an unexported field, where `SetInt` and the other setters panic.

```go
s := Stow{Station: "OST"}
v := reflect.ValueOf(&s).Elem()
r := v.FieldByName("Reason")
fmt.Println(r.CanSet()) // true
r.SetString("gale")
c := reflect.ValueOf(s).FieldByName("Reason")
fmt.Println(c.CanSet()) // false: a copy
```

## Your task

`applyDefaults` changes nothing. Pass each field that can be set, is zero and has a `default` tag to `setText`, in `defaults.go`, with the tag's text, and return any error it gives, with the field's name.

`defaults: want a pointer to a struct, got %T` is the error for a `ptr` that doesn't point to a struct.

> [!NOTE]
> Zero now always means "use the default", so tlmd's `Buffer: 0` becomes 16.
