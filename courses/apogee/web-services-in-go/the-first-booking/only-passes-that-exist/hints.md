`knownSat` and `knownStation`, in `data.go`, answer the first two rules. Check them in that order, before you look at the catalog.

Loop over `catalog` for the entry whose satellite and station equal `b`'s and whose `AOS` and `LOS` are `Equal` to `b`'s. If there is one, its bookings closed ten minutes before its AOS: compare `now` with that moment. If there is none, it is the third message.

If the booking in Furrowlight's time isn't found, times are compared with `==` somewhere. If the pass is still open at the very second bookings close, your comparison counts that second as open, but from that instant the pass is closed. Write each time in UTC with the layout `time.RFC3339`.
