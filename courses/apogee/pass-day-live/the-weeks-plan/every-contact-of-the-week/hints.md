Read every input before you print anything. A `bufio.Scanner` gives the archive a line at a time; count the lines from 1 for the error, and `strings.Fields` splits a line into its four fields. The pass ID's first two parts, split at `/`, are the satellite and the station.

A candidate is at a rehearsed station, `canServe` says yes for its satellite (find it in `newNetwork().Satellites`) and its maximum elevation is at least 10. Give `assign` one station's candidates at a time, sorted with `byAOS`, through `slices.Values`, with `nil` for the dishes' free times and two minutes of turnaround.

Put every station's contacts in one slice and sort it with `byAOS` before you write, or the lines come out station by station. A day's file holds the contacts whose AOS is from its 00:00 to the next day's, with times as `15:04:05`; make `plan/` with `os.MkdirAll`, write each file with `os.CreateTemp` in the same folder, then `os.Rename` it into place.
