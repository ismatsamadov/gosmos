Build the formula in pieces: convert `elevDeg` to radians first, then give R+h and R·cos e names of their own before you combine them.

`math.Cos` and `math.Sin` take the radians from `radians(elevDeg)`, and `math.Sqrt` takes the difference of the two squares, each a value times itself.

Skeleton: `e := radians(elevDeg)`, `a := ___ + altKm`, `b := ___ * math.Cos(e)`, then `return math.Sqrt(___) - ___ * math.Sin(e)`.
