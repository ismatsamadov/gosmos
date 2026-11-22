---
slug: quick-check-bodies-and-chunks
title: Quick Check: Bodies and Chunks
type: quiz
position: 2
difficulty: 8
topics: http, streams
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Monday evening. Nkem remembers the booking service answering 413 to every body over 65,536 bytes, and before the recorder streams its first real pass, she asks from Port Ellery where a body ends. Go's client sends 24 bytes chunked like this:

```text
c\r\nfirst piece \r\n
c\r\nsecond piece\r\n
0\r\n\r\n
```

## Your task

Answer Nkem's three questions: where a `Content-Length` body ends, where a chunked one ends, and why the recorder sends its passes chunked.
