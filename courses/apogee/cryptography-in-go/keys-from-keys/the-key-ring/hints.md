Read the lines with a `bufio.Scanner`, counting them from 1 as you go, and skip a line that is empty or starts with `#` before you hand it to `parseVersion`. Then give the version to `ring.add`.

Wrap an error from either call so its message starts with `line `, the number and a colon, keeping the original with `%w`, and return a nil ring with it. After the loop, the scanner's `Err` may hold a read error.

For `SatKey`, look through `r.versions` for the version whose `ID` is `keyID` and pass its `Key` to `satKey`. When none matches, wrap `ErrUnknownKey` with `%w` and the quoted ID, so `errors.Is` still finds it.
