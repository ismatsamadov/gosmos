Build the file's path from `cfg.Dir` and the name `bookings.json` with `filepath.Join`, and open it with `OpenFileStore`. When that fails, return a nil `*Service` and its error.

The `Store` field takes any `Store`, and a `*FileStore` is one, so it goes where `NewMemStore()` was.

`GET /bookings` is a pattern of its own beside `POST /bookings`: register it with `mux.HandleFunc` and the method value `s.listBookings`. Until you do, a GET of `/bookings` answers 405, because only POST is registered there.
