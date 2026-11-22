Which of your functions already knows d? Call it first and keep its result in a variable.

`math.Log10(x)` returns log₁₀ of x. The formula uses it twice, once on the distance and once on `freqMHz`, each time multiplied by 20.

Skeleton: `d := slantRangeKm(___, ___)`, then `return 20*math.Log10(___) + 20*math.Log10(___) + 32.44`.
