---
slug: reading-a-tag
title: Reading a Tag
type: code
position: 1
difficulty: 8
topics: reflection, strings, structs
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T03:14:43Z
---
Tuesday 29 May. Nkem has written each field's rules beside it, in a `check` tag, and the validator's next job is to read them back.

A struct tag is the text in back quotes after a field's type: `key:"value"` pairs with spaces between them. A `StructField`'s `Tag` reads a key two ways: `Get` returns its value, or "" both when the key is missing and when it's there but empty, and `Lookup` also says whether the key is there.

```go
tag := reflect.StructTag(`json:"" check:"min=1"`)
v, ok := tag.Lookup("json")
fmt.Printf("%q %v\n", v, ok) // "" true
v, ok = tag.Lookup("default")
fmt.Printf("%q %v\n", v, ok) // "" false
```

## Your task

`tagLines` writes each field's whole tag. Write `<name>: <rules>` with the field's `check` rules instead, or `<name>: no check` when it has no `check` key. `main` prints `PipelineConfig`, tlmd's ingest settings, then a `Note`, which has only a `json` tag:

```text
Workers: min=1,max=32
Buffer: min=0,max=1024
Feeds: required,len=3
Text: no check
```

> [!TIP]
> `go vet` warns about a tag that isn't in `key:"value"` form, such as `check: "min=1"`, which `Lookup` can't read.
