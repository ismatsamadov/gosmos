---
slug: no-implements-keyword
title: No implements Keyword
type: code
position: 1
difficulty: 5
topics: compiler-errors, interfaces, methods
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Paz has fitted a test radio for another frequency range, Ka-band, to La Grande, Salar Quieto's big dish. Tomás wants it in tonight's estimate without changing `downlinkSeconds`.

Go has no `implements` keyword. A type satisfies an interface as soon as it has every method the interface lists, and until then the compiler won't use it as one:

```go
type Namer interface{ Name() string }

type Dish struct{ Nickname string }

func (d Dish) Name() string { return d.Nickname }

// Dish has Name, so it is a Namer.
var n Namer = Dish{Nickname: "Chiquita"}
```

## Your task

The program doesn't build yet. Read the compiler's message, then give `KaBandRadio` the method it's missing, with a rate of 300 Mbps. Don't change `Radio` or `downlinkSeconds`.
