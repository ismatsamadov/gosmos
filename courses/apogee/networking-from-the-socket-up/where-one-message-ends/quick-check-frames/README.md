---
slug: quick-check-frames
title: Quick Check: Frames
type: quiz
position: 7
difficulty: 7
topics: networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Friday morning, Mirela reads dishd's frame reader from Port Ellery. She reviews anything that talks to the dishes, and opens with her usual line: "Your code reads the frame. Now show me the one it refuses." Then three questions about this:

```go
br := bufio.NewReader(conn)
name, err := readBanner(br)
if err != nil {
	return err
}
payload, err := readPayload(br)
```

## Your task

Answer Mirela's three questions: where the first frame's bytes can be once `readBanner` returns, which reader must read every frame after it, and when a bad length is refused.
