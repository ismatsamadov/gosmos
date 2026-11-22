---
slug: quick-check-signed-requests
title: Quick Check: Signed Requests
type: quiz
position: 3
difficulty: 8
topics: http, security
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Mirela has read the sample receiver, and before it goes to customers she has three questions about a signed event as it arrives:

```text
POST /hooks/apogee HTTP/1.1
Content-Type: application/json
Apogee-Timestamp: 1831732073
Apogee-Signature: 05ed…(64 hex digits)

{"id":"evt-000001","type":"data.landed",…}
```

## Your task

Answer three questions: why the timestamp is signed along with the body, when the receiver may trust the `booking_id` in the body, and why a time up to 300 seconds from its own clock is accepted.
