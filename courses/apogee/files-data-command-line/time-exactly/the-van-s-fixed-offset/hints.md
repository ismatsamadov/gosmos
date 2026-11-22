Make the location first, with `time.FixedZone`: an empty name is fine, and the offset is already in the seconds it wants.

Then show `t` in that location with `In`. Adding the offset to `t` would change the instant, and the test checks the instant with `Equal`.
