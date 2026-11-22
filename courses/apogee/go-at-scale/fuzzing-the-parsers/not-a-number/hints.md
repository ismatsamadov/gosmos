The range check is in `reading`, which both readings go through, so one change covers the battery and the temperature.

`if !(v >= lo && v <= hi)` writes the condition a good value meets, then refuses whatever doesn't meet it. For NaN, `>=` and `<=` are both false, so NaN lands on the refusing side.

If NaN still gets through, the condition still names the bad values with `<` and `>`, and NaN is never one of them: name the good values, then negate.
