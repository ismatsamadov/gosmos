---
slug: benchmarks-that-tell-the-truth
title: Benchmarks That Tell the Truth
position: 8
---
Write benchmarks with `for b.Loop()`, run them with `testing.Benchmark`, keep setup and input copies out of the measurement with the timer, report your own metrics and allocations, check that two versions agree before comparing them, and compare them with sub-benchmarks.
