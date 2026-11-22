---
slug: quick-check-rate-limits
title: Quick Check: Rate Limits
type: quiz
position: 2
difficulty: 8
topics: testing, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra is writing the customer guide's page on fair use, and wants it to promise exactly what the token bucket does, no more. She traces Brinewatch's morning on her notepad:

```text
09:00:00.0  10 requests: all answered
09:00:00.0  1 more: 429, Retry-After 1
09:00:02.5  how many now?
```

## Your task

Answer Yusra's three questions: what a bucket of 10 promises a customer, how many of a burst at 09:00:02.5 are answered, and why the bucket's tests never sleep.
