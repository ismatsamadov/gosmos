Make the map with `make`, sized for the bookings, and fill it in one loop over them: each booking's key is a `bookingKey` of its satellite, its station and `AOS.Unix()`.

Store an ID only when its key isn't in the map yet: the two-value lookup, `_, ok := index[k]`, says whether it is. Then loop over the passes, make each pass's key the same way, and append a `Planned` with the ID the map gives for it.

If a pass booked twice gets the second booking, the map's entry was overwritten. If the booking written in +05:30 isn't found, the key holds a `time.Time` instead of its Unix seconds.
