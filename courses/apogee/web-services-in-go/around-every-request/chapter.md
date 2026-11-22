---
slug: around-every-request
title: Around Every Request
position: 5
---
Middleware: functions that wrap an `http.Handler` to add headers, carry a request ID in the request's context, record the status through a wrapped `ResponseWriter`, write an access log, recover panics, and chain in a fixed order.
