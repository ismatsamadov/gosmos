What does `reg[id]` give back for `brine-2`, and what for `pipit-3`? The test wants different answers for the two, so the value alone can't decide.

Use the comma-ok form, `value, ok := m[key]`, and decide on `ok`: it is `false` only when the key is missing.
