Call `c.now()` once at the top and use that time throughout. Lock `c.mu` only around the map: a lookup with comma-ok, then unlock before calling `c.fetch`, which may take a while and must not keep other stations waiting.

The reading is fresh while `now.Sub(k.at)` is less than `c.ttl`; at exactly a minute it is not, which is when the cache asks again.

After a successful fetch, lock again and store a `keptReading` with the reading and `now`. When the fetch fails, return `Weather{}` and the error, and leave `c.kept` as it was.
