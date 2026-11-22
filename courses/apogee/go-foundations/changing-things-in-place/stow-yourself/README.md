---
slug: stow-yourself
title: Stow Yourself
type: code_tests
position: 4
difficulty: 5
topics: methods, pointers
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:26Z
---
Sunniva wants every stow to say why: `yst.Stow("wind 81 km/h")` should park the dishes and record the reason. Both methods run, and neither changes a thing.

A method with a value receiver gets a copy, like any argument. A method that changes its value needs a pointer receiver, `*T`:

```go
func (s *Satellite) SetMode(m Mode) {
	s.Mode = m
}

sat := Satellite{ID: "pipit-2", Mode: "NOMINAL"}
sat.SetMode("SAFE")
fmt.Println(sat.Mode) // SAFE
```

`sat` isn't a pointer, yet the call works: Go takes its address for you.

## Your task

Give `Stow` and `Unstow` pointer receivers, so that they change the station they're called on.
