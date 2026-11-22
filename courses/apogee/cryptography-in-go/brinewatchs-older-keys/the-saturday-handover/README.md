---
slug: the-saturday-handover
title: The Saturday Handover
type: code_tests
position: 7
difficulty: 10
topics: crypto, errors, security
created: 2026-11-13T03:46:40Z
modified: 2026-11-13T03:46:40Z
---
Saturday 7 October, 06:00 UTC. Brinewatch's window opens, and its security office's request arrives with the new receiving key beside it: one line naming that key by its whole fingerprint, signed by the key whose 64 digits Yusra checked by phone.

The pieces are ready in `brinewatch.go`; what's left is their order. Nothing in the request counts until its signature checks out, and the key beside it counts only if its fingerprint is the one the signed line names, so a key swapped in on the way gets no secret:

```text
brinewatch handover sha256:d249…8c1e
```

## Your task

Write `handOver`: check `r.Sig` over `r.Text` with `checkRequest`, read `r.Key` with `readRSA`, and refuse unless `r.Text` is exactly `brinewatch handover ` followed by that key's `keyFingerprint`. Then lock each secret with `encryptTo` under its label, and return two lines, or an error and no lines:

```text
api-key <hex of the locked API key>
webhook-secret <hex of the locked webhook secret>
```
