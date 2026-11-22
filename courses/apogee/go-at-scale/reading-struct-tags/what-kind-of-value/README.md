---
slug: what-kind-of-value
title: What Kind of Value
type: code
position: 0
difficulty: 8
topics: reflection, structs, types
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T09:50:21Z
---
The validator will be handed structs it has never seen, so its first skill is looking inside one, and you start it on passplan 2's settings, `Config`.

`reflect.TypeOf(v)` returns the type of the value an `any` holds. A struct type's `Fields()` ranges over its fields, each a `reflect.StructField` with a `Name` and a `Type`. A type also has a kind, what it is built from: `time.Duration` is a type of its own, and its kind is `int64`.

```go
type Band string

t := reflect.TypeOf(Band("S"))
fmt.Println(t)        // main.Band
fmt.Println(t.Kind()) // string
```

## Your task

`fieldLines` returns `nil`. Range over the fields of `v`'s type and return a line for each: its name, its type and its kind, with a space between them. `main` prints the lines for `Config`:

```text
Slots int int
Days int int
MinElevDeg float64 float64
Turnaround time.Duration int64
Replan time.Duration int64
RunLimit time.Duration int64
Queue int int
Quiet time.Duration int64
```

> [!NOTE]
> `Fields` came in Go 1.26. Older code counts the fields with `NumField` and reads each one with `Field(i)`.
