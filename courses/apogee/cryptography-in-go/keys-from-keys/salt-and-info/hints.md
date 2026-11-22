Find the call to `hkdf.Key` inside the loop. Its fourth argument is the info, and it is the same text on every pass of the loop.

Each line's name is already in `label`, and the info is exactly that text, bar and all.

If two lines still print the same bytes, the info is still a fixed string. If all four differ but not as expected, check that the master and the salt are as the starter had them.
