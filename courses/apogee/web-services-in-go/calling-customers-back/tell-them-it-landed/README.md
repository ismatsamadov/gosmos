---
slug: tell-them-it-landed
title: Tell Them It Landed
type: code_tests
position: 0
difficulty: 8
topics: http, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Furrowlight's program asks `GET /bookings` every five seconds, 17,280 times a day, to learn whether its image has arrived in full, which Apogee calls a landing. Yusra's TL;DR: "Tell them once, when it lands."

A webhook turns the asking around: when something happens, the service sends a `POST` to an address the customer chose. The receiver answers any 2xx to say it arrived: 200, 202 and 204 all count. It is built like any other request:

```go
req, err := http.NewRequestWithContext(ctx,
	http.MethodPut, u, strings.NewReader(note))
if err != nil {
	return err
}
req.Header.Set("Content-Type", "text/plain")
```

## Your task

`post` marshals `ev` but never sends it. Send it through `c` as a `POST` to `url`, made with `ctx` and carrying `Content-Type: application/json`; close the answer's body; and return an `*AnswerError` holding the status for any answer outside 200 to 299. `secret` and `ts` are for tomorrow.

> [!WARNING]
> A check for `200` alone turns away a receiver's `204 No Content`, which means the event arrived.
