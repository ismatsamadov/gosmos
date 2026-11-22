---
slug: the-record-says-how
title: The Record Says How
type: code_tests
position: 3
difficulty: 7
topics: crypto, errors, strings
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Laptops get faster every year, so Mirela plans to raise the count next year. A record holding only a salt and a key can't say which count made it, and the day the count changes, every older password would stop checking out.

So each record names everything needed to check it again: the algorithm, the count, the salt and the key, joined by `$`, the bytes in lowercase hex. A reader of such records refuses whatever it wasn't written for, rather than guess:

```go
rec := fmt.Sprintf("v%d$%x", 2, []byte("salt"))
fmt.Println(rec) // v2$73616c74
f := strings.Split(rec, "$")
fmt.Println(len(f), f[1]) // 2 73616c74
```

## Your task

The console stores each record as one line:

```text
pbkdf2-sha256$<count>$<32 hex>$<64 hex>
```

1. Make `formatRecord` write a `Record` in this form.
2. Make `parseRecord` read one back, returning an error for any other algorithm, a count outside 1,000 to 10,000,000, a salt that isn't 32 hex digits, a key that isn't 64, and a missing or extra field.

> [!NOTE]
> A record is checked with its own count, so an older one keeps working until its owner's next login replaces it.
