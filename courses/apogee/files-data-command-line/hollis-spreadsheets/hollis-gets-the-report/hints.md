Call `stationColumn` first and return its error before writing anything: a sheet without a STATION column, or with a broken quote, gets no report at all.

Count the stations in a `map[string]int`, then build the rows: the header, one row per station with its keys sorted as in the readme's example, each count turned into text with `strconv.Itoa`, and `TOTAL` last.

The total is the number of stations read, one per missed pass. Hand all the rows to `writeSheet` and return its error, so a full stick is reported too.
