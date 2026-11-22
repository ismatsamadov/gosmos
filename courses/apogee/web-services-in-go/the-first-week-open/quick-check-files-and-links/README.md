---
slug: quick-check-files-and-links
title: Quick Check: Files and Links
type: quiz
position: 4
difficulty: 9
topics: http, security
created: 2026-11-13T06:35:19Z
modified: 2026-11-13T16:15:49Z
---
Thursday afternoon, before Hollis's field team gets its first link, Nkem has three questions about the files the service now takes in and gives out:

```text
POST /bookings/{id}/commands  token, a 64 KiB form
GET  /bookings/{id}/file      token, Range
POST /bookings/{id}/link      token
GET  /files/{id}?expires&sig  the link alone
```

## Your task

Answer three questions: what a client does when a resume gets `200` instead of `206`, why `readCommands` limits `r.Body` before it reads the form, and why a link's signature covers its expiry.
