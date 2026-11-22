---
slug: quick-check-errorf-or-fatalf
title: Quick Check: Errorf or Fatalf
type: quiz
position: 3
difficulty: 4
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
After your first review, Mirela leaves two questions on this part of a test. `parsePassID` splits a pass ID like `pipit-1/PEL/04182` into its satellite, station and orbit number:

```go
id, err := parsePassID("pipit-1/PEL/04182")
if err != nil {
	t.Errorf("parsePassID: %v", err)
}
if id.Station != "PEL" {
	t.Errorf("Station = %q, want %q",
		id.Station, "PEL")
}
```

Pick one answer for each question.
