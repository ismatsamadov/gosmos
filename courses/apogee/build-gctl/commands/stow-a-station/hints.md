Check `len(args)` before you read `args[0]`, then use `stationIndex` for the code: -1 means an unknown station. Change the station through `s.stations[i]`, or a pointer to that element, never through a copy.

The reason is every word after the code, put back together by `strings.Join` with one space between words, however many were typed. With no words after the code, the reason is empty and the line reads `CWR stowed`. Add `"strings"` to the imports of `commands.go`.

Check whether the station is already stowed before changing anything. On success, print the station's code followed by `stowNote` of the changed station, then a newline.
