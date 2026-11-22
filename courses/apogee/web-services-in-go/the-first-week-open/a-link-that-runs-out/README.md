---
slug: a-link-that-runs-out
title: A Link That Runs Out
type: code_tests
position: 3
difficulty: 9
topics: http, security, time
created: 2026-11-13T06:35:18Z
modified: 2026-11-13T16:15:49Z
---
Thursday, 11:00. Hollis wants Furrowlight's field team, agronomists at the farms with laptops and no API key, to fetch Tuesday's image. A link that works without a key must open one file, stop working after a day, and let nobody change either.

`POST /bookings/{id}/link` makes such a link: the file's path, its expiry in Unix seconds, and `linkSig`'s HMAC-SHA256 of both under the service's link key, which only the service holds:

```text
/files/bk-4912?expires=1833274800&sig=2757ac…
```

Change the path or the expiry and the signature no longer matches.

## Your task

`checkLink` accepts any link. Recompute the signature with `linkSig` from the path and the expiry, compare it with `sig` in constant time with `hmac.Equal`, and return `errBadLink` when they differ or the expiry isn't a whole number. Then return `errLinkExpired` from the expiry's second on.

> [!NOTE]
> The signature comes first: only a genuine link can be expired.
