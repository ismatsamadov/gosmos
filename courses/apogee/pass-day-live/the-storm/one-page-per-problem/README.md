---
slug: one-page-per-problem
title: One Page per Problem
type: code_tests
position: 2
difficulty: 10
topics: maps, metrics, time
created: 2026-11-12T09:13:52Z
modified: 2026-11-12T11:24:16Z
---
Wednesday 18 July. Ingrid: "A storm pages whoever is on call once, not once a pass. In the investors' week, that's you."

A pager rule keeps state: it fires at one count and stays firing until the count drops below another. `rota.go` gives `onCall`, the July rota.

## Your task

Every minute of `w.Now()` until the week is done, evaluate these rules for each rehearsed station, counting failed passes (not skipped ones) by LOS:

```text
rule           LOS in    fire  clear  repeat
failed-pass    last min  1     <1     never
pass-failures  last 6 h  2     <1     24 h
```

`failed-pass` goes on the morning list, never RESOLVED.

`pass-failures` pages whoever is on call at that minute, its RESOLVED to the same person.

After the week, write `reports/alerts.txt`, a line per alert in time order, a minute's by station in `newNetwork`'s order, then rule name:

```text
Thu 10:56 FIRING failed-pass YST 1 to morning list
Thu 10:56 FIRING pass-failures YST 2 to No. 039
Thu 16:56 RESOLVED pass-failures YST 0 to No. 039
```

Log each as `alert sent` (INFO) in `logs/pager.log`, `site` `HQ`, with `rule`, `station`, `kind`, `value`, `to`. The tests check Ystmark's storm, and pages at Monday 07:59 and 08:00.
