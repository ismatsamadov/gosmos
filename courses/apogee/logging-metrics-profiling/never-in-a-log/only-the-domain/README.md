---
slug: only-the-domain
title: Only the Domain
type: code_tests
position: 5
difficulty: 7
topics: logging, security, strings
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Hollis's address turns up too. When bookingd's webhooks, its messages to a customer's programs, keep failing, it writes to the customer's contact address, and its lines about those failures carry that address under `email`. An address like Hollis's says who a person is and how to reach them: it is personal data, and it is not Apogee's to copy into every log.

A line needs only enough to tell customers apart: the first character, `***@` and the domain. `strings.Cut` splits a text at the first separator and says whether it found one:

```go
a, b, ok := strings.Cut("brine-3/LAN/14384", "/")
fmt.Println(a, b, ok)
// brine-3 LAN/14384 true
```

## Your task

`maskEmail` returns the address whole. Make it return the first character of the part before `@`, then `***@` and the domain: `hollis.vey@furrowlight.example` becomes `h***@furrowlight.example`. Anything that isn't one address, with one `@` and text on both sides, becomes `[redacted]`.

> [!NOTE]
> Every address in these lessons ends in `.example`, a domain reserved for examples, so none of them is real.
