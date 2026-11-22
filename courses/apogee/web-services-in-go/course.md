---
slug: web-services-in-go
title: Web Services in Go
kind: Course
position: 3
short: Build Apogee's booking service with net/http, from a first request to a launch that drops nothing.
---
Customers book passes by email, and Yusra copies every booking into a spreadsheet by hand. On 1 February Apogee opens a booking service that customers' programs call day and night, and you build it, with Nkem reviewing. You learn HTTP from both ends: clients that never wait forever, handlers and routes, JSON checked field by field, middleware that gives every request an ID and a line in the log, bookings that survive a restart, keys and signed tokens, fair rate limits, webhooks that sign and retry, and a shutdown that drops nothing in the middle of a live pass. Standard library only, and every server runs on your own Mac.
