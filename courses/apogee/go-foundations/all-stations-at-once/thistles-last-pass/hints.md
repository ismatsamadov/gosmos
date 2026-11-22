With `merge` returning `feeds[0]`, only Port Ellery's beacons reach the map. What does `merge` need for each feed, and when may it close its output?

`merge` needs a `sync.WaitGroup`, one forwarding goroutine per feed started with `wg.Go`, and one more goroutine that waits for them and closes the output. For the order, `slices.Sort` sorts a slice of strings in place.

If the output is right on one run and wrong on the next, the station lists are still in arrival order: sort `heardBy[seq]` before it is joined.
