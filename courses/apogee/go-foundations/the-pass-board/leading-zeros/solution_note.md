`%02d` prints a whole number at least two digits wide and fills the gap with zeros, so 4 becomes `04` while 23 stays `23`. `Sprintf` returns the result, which is what `clock` has to hand back.
