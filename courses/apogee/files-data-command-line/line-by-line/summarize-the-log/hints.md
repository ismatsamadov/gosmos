Make the `ByDish` map before the loop, and count each line at the top of the loop, before `isReading` can skip it.

For each reading, raise `Readings` and the dish's entry in `ByDish`. Replace `MaxElevDeg` when this is the first reading or its elevation is higher than the one kept.

A log whose readings are all below zero must report the highest of them, not 0. If a test sees 0 there, make sure the first reading always sets `MaxElevDeg`, whatever its value.
