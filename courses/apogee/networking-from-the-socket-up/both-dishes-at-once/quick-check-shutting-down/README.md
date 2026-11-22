---
slug: quick-check-shutting-down
title: Quick Check: Shutting Down
type: quiz
position: 6
difficulty: 8
topics: goroutines, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Friday 28 April, the morning of the deploy. Nkem's review comes from Port Ellery as three questions about this first draft of `Shutdown`:

```go
func (s *Service) Shutdown() {
	s.ln.Close()
	s.wg.Wait()
	s.Watchers.CloseAll()
	for _, d := range s.Dishes {
		d.Close()
	}
}
```

## Your task

Answer Nkem's three questions: what a blocked `Accept` does when its listener closes, what makes a blocked `Read` return, and why this draft never returns while watchers are connected.
