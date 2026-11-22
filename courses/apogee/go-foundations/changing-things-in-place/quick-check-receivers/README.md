---
slug: quick-check-receivers
title: Quick Check: Receivers
type: quiz
position: 6
difficulty: 4
topics: methods, pointers
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:08:26Z
---
Nkem's review of your methods comes back as two questions instead of comments; she'd rather you find the answers than read hers. Both are about this code:

```go
func (p Pass) Duration() time.Duration {
	return p.LOS.Sub(p.AOS)
}
func (s *Station) Stow(reason string) {
	s.Stowed, s.StowReason = true, reason
}
cwr := Station{Code: "CWR", Name: "Cape Wren"}
cwr.Stow("wind 75 km/h")
```

## Your task

Pick one answer for each question.
