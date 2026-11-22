---
slug: the-night-of-eleven-pages
title: The Night of Eleven Pages
type: code_tests
position: 7
difficulty: 9
topics: metrics, structs, time
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Friday 5 January. The new rules go live on Monday at 08:00, and first Ingrid replays the night of 27–28 December through them: `replay`, in `replay.go`, gives the pager each rule's value at Lantern Atoll every minute from 19:00 to 09:00. "Eleven pages that night. I want one."

Every idea of the week is now a field of a `Rule`: `FireAt` and `For` hold a threshold for a time, `ClearBelow` is the lower line, `Repeat` the reminder and `Severity` the reader, while the window is in the metric the replay counts. Ingrid's table:

```text
rule           fire  clear  for  repeat  severity
pass-failures  2     1      0    24h     Critical
failed-pass    1     1      0    0       Warning
disk-full      0.90  0.85   15m  12h     Critical
```

## Your task

`newRules` still returns the old rules' settings. Return Ingrid's three rules instead, and the replay sends exactly eight alerts: one page to Tomás at 20:21, its RESOLVED at 08:45, and six lines for the morning list.

> [!NOTE]
> The replay gives `pass-failures` the failed passes of the last six hours, `failed-pass` those of the last minute, and `disk-full` the share of the disk in use.
