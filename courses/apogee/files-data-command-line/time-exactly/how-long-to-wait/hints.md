Parse first and return any error from `time.ParseDuration` as it is. Text such as `30` without a unit already fails there.

A parsed duration can still be zero or below. Compare it with 0 and return an error of your own for those, made with `fmt.Errorf`.
