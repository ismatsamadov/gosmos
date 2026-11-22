---
slug: optional-notes
title: Optional Notes
type: code_tests
position: 5
difficulty: 4
topics: structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis runs operations at Furrowlight, the customer whose satellites photograph farmland, and he answers the van's first replies in capitals: NO EMPTY NOTES. But `"confirmed":false` must stay, because on his sheet a missing value means the van never answered.

`,omitempty` after a key name leaves the field out when its value is empty: `""` for a string, 0 for a number, `false` for a bool. That suits a field whose absence says nothing, and is wrong for one where `false` or 0 is an answer.

```go
type Stow struct {
    Code   string `json:"code"`
    Reason string `json:"reason,omitempty"`
}
// Stow{Code: "CWR"} encodes as {"code":"CWR"}
```

## Your task

1. In `Confirmation`, add `omitempty` to the note's tag only, so that `confirmed` appears on every line, false included.
2. Implement `encodeConfirmation` with `json.Marshal`.
