Check the name first. An empty one means the bookings are already arriving on `stdin`, so hand it to `readBookings` and return what it returns.

For a name, open the file with `os.Open` and return its error before reading anything. Close the file with `defer` once it is open, then read it with `readBookings`.

If the file test fails because stdin was read, the function reads stdin even when a name is given: only an empty name should touch it.
