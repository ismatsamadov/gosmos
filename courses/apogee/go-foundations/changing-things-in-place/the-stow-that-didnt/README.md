---
slug: the-stow-that-didnt
title: The Stow That Didn't
type: code
position: 0
difficulty: 4
topics: functions, pointers
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Sunniva, from Ystmark: "Storm by 03:00. Stow both dishes." Your code runs `stow(yst)`, and the Board still shows Ystmark tracking.

A function gets a copy of each argument, so `stow` parks a copy. To change the caller's variable, pass its address: `&x` is a pointer to `x`, a parameter of type `*T` receives it, and `*p` is the variable it points to.

```go
func addOrbit(count *int) {
	*count++
}

orbits := 15
addOrbit(&orbits)
fmt.Println(orbits) // 16
```

Through a pointer to a struct, `s.Stowed` reaches the field directly: Go follows the pointer for you.

## Your task

Change `stow` to take a `*Station`, and call it with `&yst`.

> [!NOTE]
> A pointer that holds no address is `nil`, the zero value of every pointer type.
