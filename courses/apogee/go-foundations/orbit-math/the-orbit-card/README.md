---
slug: the-orbit-card
title: The Orbit Card
type: code
position: 8
difficulty: 5
topics: formatting, functions
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Pipit-2 rises over Salar Quieto at 10:42 UTC. Tomás wants every number for the pass on one card, computed by your toolbox instead of copied from his notebook.

A call that returns one value can sit inside another call, in arithmetic, or among `Printf`'s arguments:

```go
fmt.Printf("Orbits per day: %.2f\n",
	1440/periodMinutes(525))
```

## Your task

Under the title in `main`, print the card's other four lines, one `Printf` each, calling the functions with Pipit-2's constants. `dopplerHz` works in Hz, but the card shows kHz:

```text
Pipit-2 over Salar Quieto
Period: 95.50 min
Slant range at 10 deg: 1815.08 km
Doppler at -6.5 km/s: +48.68 kHz
Path loss at 10 deg: 164.64 dB
```

> [!TIP]
> `%+.2f` prints a sign even when the number is positive.
