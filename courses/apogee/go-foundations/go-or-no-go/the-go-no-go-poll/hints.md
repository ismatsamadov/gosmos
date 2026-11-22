When the station is down and the wind blows at 90 km/h, which reason must the poll report? Put that check where it runs first.

A `switch` with no value, or a row of `if` statements that each `return`, stops at the first condition that is true. `!up` is true when the station is down.

Shape: `case !up:` returns the first reason, then one case each for `windKmh ___ 72`, `maxElevDeg ___ 3` and `battV ___ 6.6`, in that order, and `default:` returns `"GO"`. Check each blank against the tests that sit exactly on a limit.
