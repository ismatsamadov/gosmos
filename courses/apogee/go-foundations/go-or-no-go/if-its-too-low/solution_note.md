Each `if` checks `maxElevDeg` after that pass's value is assigned, so only Pass 2's 7.5 is under 10 and keeps its warning, while the `Printf` lines outside the braces print for every pass.
