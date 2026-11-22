---
slug: quick-check-request-lines
title: Quick Check: Request Lines
type: quiz
position: 2
difficulty: 8
topics: http
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tuesday. Last winter Go's `net/http` read request lines for the booking service, through which customers book passes, and handed each handler its method and path. Now the station server reads them itself, byte by byte, and Nkem, who reviewed the booking service too, asks three questions from Port Ellery about these three lines:

```text
GET /recordings HTTP/1.1
GET /recordings HTTP/1.0
get /recordings HTTP/1.1
```

## Your task

Answer Nkem's three questions: what the target is, what a 505 says that a 400 doesn't, and what becomes of `get`.
