---
slug: once-and-only-once
title: Once and Only Once
position: 2
---
Initialize once with `sync.Once`, `sync.OnceValues` and `sync.OnceFunc`; count with `atomic.Int64`, publish with `atomic.Pointer` and replace only with something newer by compare-and-swap; and see, through Go's memory model, why a program without a race behaves.
