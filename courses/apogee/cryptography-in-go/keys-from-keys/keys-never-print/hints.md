Give `Key` value receivers, so its methods work on a key and on a pointer to one alike. `String` returns the ID followed by ` (secret hidden)`, and `GoString` can build its text with `fmt.Sprintf`, quoting the ID with `%q`; add `"fmt"` and `"strings"` to the file's imports.

The ring's methods have pointer receivers, as `*KeyRing` says. For `String`, collect the `ID` of each of `r.versions` in order and join them with `strings.Join` and `", "` after `key ring: `; `GoString` can return what `String` returns.

If `TestRing` still finds a secret while the other tests pass, look at the verb in its message: `%#v` needs the ring's `GoString`, and `%v`, `%+v` and `%s` its `String`. Inside the unexported field `versions`, `fmt` can't reach `Key`'s methods.
