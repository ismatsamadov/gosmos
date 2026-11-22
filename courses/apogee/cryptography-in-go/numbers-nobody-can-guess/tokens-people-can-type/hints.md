`rand.Text` does all of the work: it takes no arguments and returns a new token each time you call it.

Replace the whole `return` line with one that returns `rand.Text()`, then remove the `_ = rand.Text` line and the imports nothing uses anymore.

If the build stops at `"strconv" imported and not used` or the same for `"time"`, delete those imports: the new token needs nothing from the clock.
