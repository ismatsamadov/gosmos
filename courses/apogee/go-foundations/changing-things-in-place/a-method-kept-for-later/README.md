---
slug: a-method-kept-for-later
title: A Method Kept for Later
type: code_tests
position: 5
difficulty: 5
topics: functions, methods, pointers
created: 2026-11-04T17:00:51Z
modified: 2026-11-06T16:51:26Z
---
Sunniva will be outside checking Ystmark's heaters when the wind rises. The checklist needs to keep that station's stow action now and call it later, with the reason known then.

`s.Stow` without parentheses is a method value: a function with its receiver already bound. With a pointer receiver it keeps the pointer, so calling it changes the same station. A value receiver would keep a copy of its receiver instead.

```go
stamp := pass.Duration
fmt.Println(stamp()) // receiver already bound
```

## Your task

Complete `stowAction` by returning the supplied station's method value. Do not call `Stow` yet or write a new wrapper function. The caller supplies the reason when it runs the returned action.
