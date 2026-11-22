Sort a copy of `d.Slots` by AOS with `slices.SortFunc`. For a slot, the certificate in force is the station's `New` one once AOS is no longer before its `switchTime` (from the day's `Passes()`, starting at `d.SwitchFrom`), else its `Old` one; the list is `d.List` once AOS is no longer before `d.ListAt`, else nil.

`Connect` comes first, and an error there misses the pass. Then each `SealAt` time gets its own version from `signingKey`, over the ring's `Versions()`, and a `Ping` at that time plus 5 seconds; the first error misses the pass.

A line is AOS as `15:04:05`, the satellite and the station, then `flown`, the key IDs joined by `+` and the certificate's `Leaf.SerialNumber`, or `missed:` and the error. If the midnight pass shows `k-2028-07` twice, its keys were picked at AOS instead of at each seal time.
