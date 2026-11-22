Add the field the way `passes` was added, and set it in `newShell`. Codes are three capital letters, so `stationIndex(s.stations, "pel")` returns -1 too: only the table decides what a real code is.

Check the arguments first, before sorting or printing anything. More than one argument returns the second error at once; one argument that the table doesn't hold returns the first, with the code quoted by `%q`.

In the printing loop, skip passes at other stations when a code was given, and print every pass when none was. A known station with no passes tonight, such as `SQT`, prints no lines at all.
