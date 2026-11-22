---
slug: quick-check-comparing-secrets
title: Quick Check: Comparing Secrets
type: quiz
position: 5
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Before the night shift's phrase goes back into use, Mirela lists the week's three comparisons on the whiteboard and asks which of them must take the same time whatever the bytes:

```go
hmac.Equal(gotTag, wantTag)  // a command's tag
got == want                  // two file digests
subtle.ConstantTimeCompare(a[:], b[:]) // phrases
```

## Your task

Answer Mirela's three questions: why a tag is compared with `hmac.Equal`, whether two public file digests need the same care, and why phrases are hashed before they are compared.
