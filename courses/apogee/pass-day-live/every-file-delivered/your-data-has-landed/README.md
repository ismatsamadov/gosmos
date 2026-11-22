---
slug: your-data-has-landed
title: Your Data Has Landed
type: code_tests
position: 1
difficulty: 9
topics: crypto, http, json
created: 2026-11-12T05:36:26Z
modified: 2026-11-12T11:16:23Z
---
Thursday afternoon. Yusra: "TL;DR: the investors will ask Furrowlight how soon its images arrive." And Joss has Kessick's receiver waiting for Thistle-2's first science file.

Each landing is told to its customer by webhook: a POST of a `data.landed` event, the `Event` in `types.go` as `json.Marshal` writes it, to the customer whose `name` in `customers.json` is the satellite's operator. The week's first, split here:

```text
{"id":"evt-000001","type":"data.landed",
"booking_id":"bk-8001","sat_id":"cirrine-6",
"station":"ILV",
"file":"cirrine-6-20290723-0003.dat",
"bytes":6144,
"landed_at":"2029-07-23T00:04:17Z"}
```

## Your task

Read `customers.json` strictly, as you read `rehearsal.json`; its errors read like this one:

```text
customers.json: more data after the customers
```

As each booked file lands, POST its event to the customer's `webhook_url` through `w.HTTP`, IDs from `evt-000001` in landing order, `landed_at` truncated to the second, with these headers:

```text
Content-Type: application/json
Apogee-Timestamp: 1879459457
Apogee-Signature: bfeca887…6127
```

The timestamp is the attempt's Unix second on `w.Now()`; the signature is the HMAC-SHA256, in lowercase hex, of the timestamp, a dot and the body, under the customer's `webhook_secret`. Print `DELIVERED <file>` when the answer is a 2xx, and `UNDELIVERED <file>` otherwise. The tests' receivers check every request and the week's events, Kessick's for `bk-8792` among them.
