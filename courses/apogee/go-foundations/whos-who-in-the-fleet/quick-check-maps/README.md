---
slug: quick-check-maps
title: Quick Check: Maps
type: quiz
position: 3
difficulty: 3
topics: maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-22T10:08:26Z
---
The registry is ready for review. Nkem reads your map code line by line, then asks three questions about this snippet, one at a time, the way she reviews everything:

```go
seen := map[string]int{"pipit-1": 3, "brine-3": 1}
fmt.Println(seen["thistle"])
for id := range seen {
	fmt.Println(id)
}
var pending map[string]int
pending["thistle"] = 1
```

## Your task

Each question has exactly one right answer. Before you pick, think about what each line does when it runs.
