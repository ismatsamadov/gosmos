---
slug: quick-check-byte-order
title: Quick Check: Byte Order
type: quiz
position: 2
difficulty: 7
topics: bytes, conversions
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tuesday 28 March. Tomás, off the plane and still in his parka, takes a coffee in the hut and asks three questions about five bytes of Monday's status frame. "Get the order wrong and the dish points at the sea!"

```go
b := []byte{0x47, 0x95, 0x02, 0xff, 0x83}
az := binary.BigEndian.Uint16(b[0:2])
state := b[2]
temp := int16(binary.BigEndian.Uint16(b[3:5]))
```

## Your task

Answer Tomás's three questions: what `az` holds, why `state` needs no byte order, and what `temp` holds.
