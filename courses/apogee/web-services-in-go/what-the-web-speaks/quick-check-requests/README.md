---
slug: quick-check-requests
title: Quick Check: Requests
type: quiz
position: 3
difficulty: 4
topics: http
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Next week the engineers of Cirrine, the customer whose satellite measures the weather, try Nkem's stand-in: a few lines of Go that answer like the booking service will. First, Nkem checks your notes against one request:

```text
GET /passes?sat=furrow-a&station=SQT HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Accept: application/json
```

## Your task

Answer three questions about this request: what its first line holds, which method would cancel a booking, and what follows the `?`.
