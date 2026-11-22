---
slug: quick-check-retries
title: Quick Check: Retries
type: quiz
position: 7
difficulty: 8
topics: errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
On Friday Furrowlight's first webhook goes out. Before it does, Nkem has three questions about the schedule every event now follows:

```text
attempt 1 fails    wait 1s to 2s
attempt 2 fails    wait 2s to 4s
attempt 3 fails    wait 4s to 8s
...
attempt 7 fails    wait 64s to 128s
attempt 8 fails    undelivered
```

## Your task

Answer three questions: what the random part of each wait is for, what the sender does with a `404`, and why a receiver must be ready for the same event twice.
