---
slug: quick-check-never-twice
title: Quick Check: Never Twice
type: quiz
position: 6
difficulty: 8
topics: networking
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T16:42:55Z
---
Thursday afternoon. Last year Mirela, Apogee's security engineer, stopped an old signed command from going up to Cirrine-1 a second time, a replay. Now she reads the new dishd closely and asks three questions about an exchange like Tuesday's:

```text
dishd: JOG seq 122, A +0.35
       (the line drops: no answer)
dishd: HELLO
YST-2: WELCOME, last_seq 122
```

## Your task

Answer Mirela's three questions: what the missing answer says, why a command sent again keeps its number, and which command must never run twice.
