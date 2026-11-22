---
slug: the-rehearsals-settings
title: The Rehearsal's Settings
type: code_tests
position: 0
difficulty: 8
topics: errors, files, json
created: 2026-11-10T20:36:24Z
modified: 2026-11-12T02:35:45Z
---
Monday 2 July. From 23 July the investors watch the network for a week; first the team flies that week in a computer: the rehearsal, which you start today. Ingrid wants each run to say first what it flies.

The rehearsal is one function, `rehearse(ctx, w World)`; `world.go` says what `w` brings; every other name is yours. Its folder, `w.Dir`, holds `rehearsal.json`:

```json
{"week":"2029-07-23","days":7,
 "stations":["YST","TIR"]}
```

## Your task

Read `rehearsal.json` strictly: `days` is 1 to 7, 7 when absent, and no `stations` means all twenty. Print one line on `w.Out`, in one write, and return nil:

```text
rehearsal: 2029-07-23, 7 days, 20 stations
```

With one, it's `1 day` or `1 station`.

For a bad file, print nothing and return `rehearsal.json: ` and one of these; a missing file's error must still match `fs.ErrNotExist`.

```text
json: unknown field "dayz"
more data after the settings
week is required
week "2029-7-23" is not a date like 2006-01-02
days 8 is not from 1 to 7
unknown station "YSX"
```

> [!NOTE]
> The tests run it in the simulated world, Nkem's program playing every station and satellite. Their names begin with `sim` or `test`, so none clashes with yours.
