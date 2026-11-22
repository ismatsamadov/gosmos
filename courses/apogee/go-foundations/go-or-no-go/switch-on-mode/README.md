---
slug: switch-on-mode
title: Switch on Mode
type: code
position: 5
difficulty: 3
topics: switch
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
What ops does during a pass depends on the satellite's mode. NOMINAL means normal operation, so the pass can proceed; SAFE means a minimal state after a fault, so the team schedules a recovery pass to restore it. DETUMBLE (still stopping its spin) and LEOP (its first days in orbit) get beacon only; any other mode means hold.

`switch` compares one value with each `case` from the top and runs the first match; `default` runs when nothing matches. A case can list several values, separated by commas. After the matching case, the switch is done, unless that case ends with `fallthrough`:

```go
switch station {
case "YST", "CWR":
	fmt.Println("check the wind first")
case "LAN":
	fmt.Println("check the rain first")
default:
	fmt.Println("no weather check")
}
```

## Your task

Replace the `return` in `modeAction` with a `switch` on `mode` that returns each mode's action, with `DETUMBLE` and `LEOP` sharing one case, so the program prints:

```text
NOMINAL: proceed
SAFE: recovery pass
DETUMBLE: beacon only
LEOP: beacon only
REBOOTING: hold
```
