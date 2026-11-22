Inside `average`, `readings` is a slice of `float64`. What do you already know how to do with a slice of numbers?

Add the readings up with a `range` loop, then divide by `len(readings)`, converted to `float64`.

Handle `len(readings) == 0` first, with an early `return`, before you ever divide.
