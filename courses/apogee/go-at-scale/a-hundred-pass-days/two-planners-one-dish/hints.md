Start from an empty `map[string][]time.Time` and range over `cur.Contacts`, breaking out of the loop at the first contact whose AOS isn't before `at`.

For each contact, take its station's slice, append zero times until it has `c.Dish` entries, and set the dish's entry to the contact's LOS plus `turnaround` if that is later than what it holds. Store the slice back in the map, since `append` may have made a new one.

If the 20:40:16 test fails, contacts that rise at `at` or later are being counted, Furrow-G's first among them; if a free time is two minutes early, the turnaround is missing.
