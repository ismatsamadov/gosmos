---
slug: quick-check-certificates
title: Quick Check: Certificates
type: quiz
position: 3
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday afternoon, with the root in the safe, Nkem draws on the whiteboard the chain every station's certificate will hang from, each one signed by the one above it. Between the stations and the root stands the stations CA, the intermediate:

```text
yst.apogee.internal
  signed by Apogee Stations CA 2028
    signed by Apogee Root CA 2028
```

## Your task

Answer Nkem's three questions: what a certificate binds together, why a certificate copying every one of Ystmark's names still fails, and why the root's private key never leaves the safe.
