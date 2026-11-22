---
slug: quick-check-status-and-bodies
title: Quick Check: Status and Bodies
type: quiz
position: 5
difficulty: 5
topics: defer, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Tonight the poll's new client meets a real night of passes. Before it does, Ingrid runs her own check on this part of it: "Three questions. One answer each."

```go
resp, err := c.HTTP.Do(req)
if err != nil {
	return Weather{}, err
}
defer resp.Body.Close()
if resp.StatusCode != http.StatusOK {
	return Weather{}, errors.New(resp.Status)
}
```

## Your task

Answer the three questions: what `Do` returns as its error for a 503, whether a body nobody reads must still be closed, and when the deferred `Close` runs.
