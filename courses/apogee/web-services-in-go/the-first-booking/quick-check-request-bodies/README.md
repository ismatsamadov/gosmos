---
slug: quick-check-request-bodies
title: Quick Check: Request Bodies
type: quiz
position: 3
difficulty: 6
topics: errors, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Mirela, who will review the whole service in January, reads the decoder before Furrowlight's first real booking. "Now show me what it refuses." She points at these lines, where `apiErr` is an `*APIError`:

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
err := dec.Decode(&b)
extra := dec.Decode(&struct{}{})

wrapped := fmt.Errorf("booking: %w", apiErr)
```

## Your task

Answer Mirela's three questions: what a misspelled field does, why the code decodes twice, and what `errors.AsType` finds in `wrapped`.
