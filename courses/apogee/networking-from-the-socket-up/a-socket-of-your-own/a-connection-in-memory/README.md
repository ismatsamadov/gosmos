---
slug: a-connection-in-memory
title: A Connection in Memory
type: code_tests
position: 5
difficulty: 6
topics: networking, testing
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Thursday. Nkem wants dishd's tests to run on any laptop without a single port, and asks from Port Ellery for one helper that every test can use.

`net.Pipe()` returns two `net.Conn` ends joined in memory: what one end writes, the other reads. A pipe has no buffer, so a `Write` waits until the other end reads it. In a test, `t.Cleanup(f)` runs `f` when the test ends, even after a failure, while a `defer` in a helper runs as soon as the helper returns:

```go
func listen(t testing.TB) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln
}
```

## Your task

`pipePair` in `pipe.go` returns `nil, nil`. Make it return the two ends of a `net.Pipe`, and close both of them when the test that asked for them ends.
