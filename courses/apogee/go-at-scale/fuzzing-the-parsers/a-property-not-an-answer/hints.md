Keep what `parse` returns: `f, err := parse(line)`. When `err` isn't nil, return at once: a refused line needs no more checks.

Give each property an `if` and a `t.Errorf`: an empty `f.SatID`; `f.BattV` that is NaN, below 0 or above 20; `f.TempC` the same with −100 and 150. For the mode, a `switch` whose one `case` lists the four modes, and whose `default` fails.

If the mutant that accepts `batt_v=nan` survives, the battery check only compares with 0 and 20, and a NaN passes both comparisons: test `math.IsNaN(f.BattV)` as well.
