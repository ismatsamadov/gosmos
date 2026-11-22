---
slug: quick-check-signatures
title: Quick Check: Signatures
type: quiz
position: 2
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Tuesday. Yusra is adding a page on receipts to her pages for customers' engineers, the customer guide, and wants every sentence exactly right before Hollis reads it. The page's sample check is two lines:

```go
pub := publishedKey("r-2028-08") // 32 bytes
ok := ed25519.Verify(pub, receipt, sig)
```

Her message opens, as always, with a summary: "TL;DR: three questions, then the page goes out."

## Your task

Answer Yusra's three questions: what a customer can do with Apogee's public key, what signing the same receipt twice gives, and why a secret both sides hold can't settle Hollis's complaint.
