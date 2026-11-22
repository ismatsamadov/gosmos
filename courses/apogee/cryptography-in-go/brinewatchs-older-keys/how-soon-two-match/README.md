---
slug: how-soon-two-match
title: How Soon Two Match
type: code
position: 1
difficulty: 8
topics: crypto, maps, printing
created: 2026-11-13T03:46:39Z
modified: 2026-11-13T12:17:19Z
---
Brinewatch's security office will read its keys' fingerprints to Yusra by phone, so nobody can swap a key on the way, and its tool shows only the first six hex digits: 3 bytes, 16.7 million codes. Tomás, at the whiteboard: "Twenty-three people in a room, and two probably share a birthday. How many keys before two share a code?"

Matching one given code takes about as many tries as there are codes. Any two match far sooner, since each new try can match every earlier one: about when the number of tries, squared, reaches the number of codes, a limit called the birthday bound.

```go
codes := math.Pow(2, 24) // six hex digits
fmt.Println(math.Sqrt(codes)) // 4096
```

## Your task

Real keys take time to make, so the texts `brinewatch 0`, `brinewatch 1`, … stand in for them. `firstMatch` compares each text's code with the first text's alone. Make it keep every code it has seen in a map, with the number of the text that had it, and return the first two texts whose codes match.

> [!NOTE]
> A whole SHA-256 fingerprint is 32 bytes: finding any two that match would take about 2^128^ tries.
