Take one line at a time: which function gives each number, and which of the constants above the functions does it take?

Call each function among `Printf`'s arguments, with `%.2f` in the format. `dopplerHz` takes and returns hertz, so multiply `pipit2DownlinkMHz` by `1e6` going in and divide the result by 1000 coming out.

The Doppler line needs its plus sign: the `+` flag, as in `%+.2f`, prints a sign whether the shift is positive or negative.
