Change the map's key type from `string` to `passKey`, and build each pass's key as a `passKey` literal instead of with `fmt.Sprintf`.

The AOS goes in as a number: `p.AOS.Unix()` gives its seconds. The rest of the loop, the lookup, the insert and the `append`, stays as it is.

If two passes that rise in the same second come out as one, the key is missing the satellite or the station: all three parts tell passes apart.
