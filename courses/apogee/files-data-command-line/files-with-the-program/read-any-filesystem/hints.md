`fs.ReadFile` takes the filesystem first and the file's name second. The name is `stations.csv`: no folder, no leading `/`.

Check the read error before parsing. Then return what `parseStations` returns: it already gives both values the function needs.

If the one-station table comes back with five stations, the function is still reading gctl's own table: read from the `fsys` it is given, not from `resources` or the disk.
