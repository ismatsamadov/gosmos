`benchInputs` returns everything `planDay` needs, in one value: call it once, at the top of `benchPlan`, before the loop.

Inside `for b.Loop() { ... }` make exactly one call: `planDay` with the six fields of the inputs, in the order its signature lists them.

If the test says `benchInputs` was called more than once, it is inside the loop, or the loop counts with `b.N`, which makes `testing.Benchmark` call `benchPlan` again and again. If it planned no times, the call to `planDay` is missing from the loop.
