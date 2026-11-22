---
slug: one-key-one-job
title: One Key, One Job
type: code_tests
position: 3
difficulty: 7
topics: crypto, errors
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Cirrine-1 will hold two keys: one that checks its signed commands, and one that opens the locked ones it gets from 25 August. Mirela: "One key, one job. A key that does two jobs is only as safe as the weaker job."

So the purpose goes into the info, beside the satellite: `sign|cirrine-1` and `seal|cirrine-1` are unrelated keys from the same master, and whatever one job might ever reveal about its key tells nobody anything about the other's. The derivation also refuses any purpose or satellite it doesn't know, so a typo never makes a key that no satellite holds:

```go
ids := []string{"YST", "PEL", "SQT"}
fmt.Println(slices.Contains(ids, "PEL")) // true
fmt.Println(slices.Contains(ids, "pel")) // false
```

## Your task

Make `satKey` derive with the info `<purpose>|<satellite>`, and return an error for any purpose but `sign` and `seal`, and for a satellite that isn't in `fleet`.
