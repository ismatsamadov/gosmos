Only the line inside `budget` changes, and no `float64` is left in it: `booked`, `10000` and `int(t)` are all whole numbers already.

Keep the order of the operations: the product of `booked` and `10000 - int(t)` comes first, and the division by 10000 last, so nothing is rounded away until the end.

If only the ten-thousand-pass case fails, with 4 instead of 5, a float is still in the arithmetic. If Port Ellery's budget is 0, the division came before the multiplication.
