---
slug: one-bug-no-outage
title: One Bug, No Outage
type: code_tests
position: 7
difficulty: 8
topics: functions, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Friday. Nkem held back the fix for Monday's bug in `getPass`, so the middleware meets the real thing: Cirrine sends the very request that broke, `GET /passes/cirrine-1/YST/1136`. This time it should get an answer, an ID, and a line in the log that Yusra can quote.

Order decides what each middleware sees: the request as the ones outside it have changed it, and the answer as the ones inside it have written it.

```go
chain(mux, requestID, logRequests(l, now))
// logRequests runs second: it finds the ID
chain(mux, logRequests(l, now), requestID)
// logRequests runs first: no ID yet
```

## Your task

`handler` chains `recoverPanics`, `logRequests`, `requestID` and `noStore` in that order, which loses the IDs from the log and the panicking request's line from it. Put them in the order that gives every request an ID, an answer and a line in the log. Monday's request should then leave two lines, its panic and its 500, both starting with `req-000002`.
