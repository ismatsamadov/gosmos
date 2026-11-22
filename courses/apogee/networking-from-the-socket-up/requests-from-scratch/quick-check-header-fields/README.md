---
slug: quick-check-header-fields
title: Quick Check: Header Fields
type: quiz
position: 6
difficulty: 8
topics: http, security
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T16:42:55Z
---
Wednesday. Mirela, at Port Ellery, has finished reviewing `readHeader`, and before she signs it off she asks three questions about five header lines from Tuesday's tests in the hut, the last one sent on purpose by a crooked client of her own:

```text
Content-Type: Text/Plain
content-type: text/plain
X-Pass: furrow-a/YST/18218
X-Pass: cirrine-1/YST/13094
Host : 127.0.0.1
```

## Your task

Answer Mirela's three questions: how names and values treat letter case, what a repeated field keeps, and why `Host :` is refused rather than read as Host.
