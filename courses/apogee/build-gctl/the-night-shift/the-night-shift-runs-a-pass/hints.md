Both changes ask the same question of a pass's station: is it stowed? A small method on `Shell` can answer it with `stationIndex` and the station's `Stowed` field.

In `cmdNext`, a pass counts only when its AOS is after now and its station is not stowed. In `cmdPasses`, print the line without its newline, then ` (stowed)` if the station is stowed, then the newline.

When a transcript fails, the test names the first line that differs and the `.golden` file it comes from: open that file in `testdata` and compare it with your output from that line on. Lines that start with `@` only move the clock and never reach `Run`.
