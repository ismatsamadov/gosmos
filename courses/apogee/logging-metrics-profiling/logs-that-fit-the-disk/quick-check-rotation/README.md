---
slug: quick-check-rotation
title: Quick Check: Rotation
type: quiz
position: 3
difficulty: 7
topics: files, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Thursday 7 December. Before the rotator goes to Lantern Atoll, Nkem asks three questions.

The first imagines the night shift's search reading `uplinkd.log`, the file open, at the moment the rotator renames it:

```text
uplinkd.log  ->  uplinkd-2028-12-07.1.log
(fresh)      ->  uplinkd.log
```

## Your task

Answer Nkem's three questions: what a reader holding the file open sees after the rename, why a rotation never cuts a record in two, and why the current file always keeps the name `uplinkd.log`.
