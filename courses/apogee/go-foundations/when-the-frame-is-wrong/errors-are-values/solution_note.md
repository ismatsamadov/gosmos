`checkBattery` returns `nil` for a healthy reading and an error made by `errors.New` for one under `minBattV`, so a caller tells the two apart by comparing the result with `nil`.
