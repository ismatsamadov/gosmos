---
slug: read-the-booking
title: Read the Booking
type: code_tests
position: 0
difficulty: 6
topics: http, json, streams
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Hollis wants his next booking to arrive as JSON, not as an email in capitals: on Friday, Furrowlight's program will send it to `POST /bookings`.

A request's body is a stream, like an open file, and `json.NewDecoder(r.Body).Decode(&b)` reads one JSON value from it into `b`. When the body isn't JSON, an empty body included, `Decode` returns an error and there is nothing worth keeping:

```go
var st Station
err := json.NewDecoder(strings.NewReader(
	`{"code":"LAN","mask_deg":8}`)).Decode(&st)
fmt.Println(st.Code, st.MaskDeg, err)
// LAN 8 <nil>
```

## Your task

`createBooking` answers 201 with an empty booking. Decode the body into `b`, and when that fails, answer 400 with `http.Error` and `request body: ` followed by the decoder's error. Then keep `b` with `s.Store.Create(r.Context(), b)`, which returns it with its new ID, and answer 201 with what it returns.

> [!NOTE]
> For any error from `Create`, answer 500 `internal error` for now: telling a double booking apart comes on Friday.
