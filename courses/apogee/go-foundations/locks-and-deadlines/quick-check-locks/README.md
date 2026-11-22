---
slug: quick-check-locks
title: Quick Check: Locks
type: quiz
position: 2
difficulty: 6
topics: defer, mutexes
created: 2026-11-01T07:36:47Z
modified: 2026-11-22T10:08:26Z
---
At your on-call handover, Mirela pulls up `Uplink`. Its only field that matters here is `mu sync.Mutex`, and the service makes one with `var up Uplink`; there is no constructor:

```go
func (u *Uplink) Send(cmd Command) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if cmd.Sig == "" {
		return ErrUnsigned
	}
	return u.transmit(cmd)
}
```

## Your task

Answer Mirela's two questions before she signs off on it.
