Decode `bookings.json` strictly into a `[]Booking`, as you did the settings; when opening it fails with `fs.ErrNotExist`, there are no bookings. A map keyed by satellite, station and AOS in Unix seconds finds a booking's contact in one lookup, and the same key into the archive tells `no such pass` from the rest.

Leave out a booking whose AOS is before the window's start or not before its end, and one whose station is in `newNetwork().Sites` but not among the stations you fly. Everything else counts, matched or not.

Try the reasons in the readme's order and stop at the first that fits: the satellite, the station, the archived pass, the bands, the 10 degrees; whatever is left was dropped. Sort the unplanned bookings by AOS, then booking ID, write AOS as `Mon 15:04:05` in UTC, and save the report before you print the count.
