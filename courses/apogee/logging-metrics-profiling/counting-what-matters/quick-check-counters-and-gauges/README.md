---
slug: quick-check-counters-and-gauges
title: Quick Check: Counters and Gauges
type: quiz
position: 2
difficulty: 7
topics: metrics
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Monday, 16:00. All afternoon the lab bench has run replayed passes through the new counters, numbers the Board will show from tomorrow, and Nkem has two questions about them.

She points at two readings of the bench's uplinkd, five minutes apart:

```text
       passes_total  passes_running
15:55            41               1
16:00             3               0
```

## Your task

Answer Nkem's two questions about these readings: what the drop in `passes_total` tells you, and why `passes_running` is a gauge and not a counter.
