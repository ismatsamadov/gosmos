`%.3f` always prints three decimals, so `2245` gains three zeros and `437.25` one, while `%.1f` rounds `61.37` to `61.4`. `%v` prints a float as short as it can, which is why the decimals wandered.
