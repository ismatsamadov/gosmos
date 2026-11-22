---
slug: every-route-in-one-place
title: Every Route in One Place
type: code_tests
position: 7
difficulty: 7
topics: http, testing
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Friday. This afternoon Yusra sends Hollis the first email with a link to a pass, and Nkem wants the whole address book checked first: "One list of routes, one table of requests."

With every route in one `routes` method, the service's addresses read top to bottom in one place, and one table can check them all. `main_test.go` sends 19 requests: each route's usual answer, a 404 for each wildcard, a wrong method on every path and a HEAD. Each row is one request and its answer:

```go
{"DELETE", "/passes", 405, "GET, HEAD",
	"Method Not Allowed"},
{"GET", "/stations/ZZZ", 404, "",
	`no station "ZZZ"`},
```

## Your task

`routes` forgets `GET /passes` and registers `/bookings/{id}` for every method, so a DELETE reads the booking instead of canceling it. Make it register exactly these routes, with `listPasses` for the list of passes and `getBooking` and `cancelBooking` for a booking:

```text
GET    /health
GET    /stations
GET    /stations/{code}
GET    /guide/
GET    /passes
GET    /passes/{sat}/{station}/{orbit}
GET    /bookings/{id}
DELETE /bookings/{id}
```
