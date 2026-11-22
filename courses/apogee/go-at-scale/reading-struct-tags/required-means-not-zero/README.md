---
slug: required-means-not-zero
title: Required Means Not Zero
type: code_tests
position: 4
difficulty: 8
topics: reflection, slices
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T03:14:43Z
---
Thursday 31 May. The host of Ilha Calmaria (ICM), the third wave's station in the calm belt of the tropical Atlantic, has sent its form back with the station's name left blank.

A rule reads one field's `Value`. `required` means the field isn't its type's zero value, and `Value.IsZero` reports that for any type: "" for a string, 0 for a number. A slice is zero only when it's nil, though, and a list someone emptied is an empty slice that isn't nil, so look at a slice's `Len` as well.

```go
var none []string
fmt.Println(reflect.ValueOf(none).IsZero())
// true
e := reflect.ValueOf([]string{})
fmt.Println(e.IsZero(), e.Len()) // false 0
```

## Your task

`checkRule` knows no rules yet and calls every rule unknown. Handle `required`: return the problem `<name>: required` when `v` is its type's zero value or an empty slice, and `nil` otherwise. Every other rule stays unknown.

> [!NOTE]
> `name` is the field's struct and name, such as `HostConfig.Name`, so each problem says where it is.
