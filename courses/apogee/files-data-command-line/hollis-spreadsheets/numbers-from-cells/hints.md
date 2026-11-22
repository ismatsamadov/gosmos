`strconv.ParseFloat` with a bit size of 64 gives the number and an error. Store the number in the row you append, next to the ID.

The loop index already counts records from 0 for the header, which Hollis's spreadsheet calls row 1, so a record's row is its index plus one.

Wrap the parsing error with `%w` after the text `row N, ELEV: `, and return `nil` rows with it: a sheet with one bad cell is refused whole.
