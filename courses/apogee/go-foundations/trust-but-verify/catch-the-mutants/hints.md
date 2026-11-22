For each mutant, look for a `seq` and `last` where it and `seq > last` give different answers. Where do `>` and `>=` disagree?

Work out each case's `want` from the real `accept`, never from a mutant: the mutant is the one that should get it wrong.

Three neighbors of `last` do the work: `seq` equal to `last`, one more than `last`, and one less.
