Each loop iteration has its own `minimum` variable under Go 1.22 and later. The three closures therefore retain thresholds of 1, 2 and 3 even when called after the loop ends.
