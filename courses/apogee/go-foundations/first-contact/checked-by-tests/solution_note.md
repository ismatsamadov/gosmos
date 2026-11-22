The test calls `stationCode` and compares what it returns with `"PEL"`. Returning exactly that string makes `got` equal `want`, so `t.Errorf` never runs and the test passes.
