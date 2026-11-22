## Why does heapProfile collect garbage right before it writes the profile?

- [x] So the profile counts what is reachable now, not an older collection's numbers
  > A heap profile reports the heap as of the last finished collection, so a fresh one makes its numbers mean memory in use now.
- [ ] Because a heap profile can only be written while a collection runs
  > The profile can be written at any time; it reports the last finished collection, however old that is.
- [ ] So the collection frees the stuck goroutines' buffers before they are counted
  > Nothing can free those buffers: the goroutines still reach them, which is why they show up in the profile at all.

## What are the four goroutines on line 15 doing?

- [ ] Flying four passes at once, one on each worker
  > A goroutine that flies a pass is inside `fly`; these four have finished flying and are stopped at the send after it.
- [x] Waiting forever to send a result that runPass stopped waiting for at LOS
  > `done` has no room, so the send waits for a receiver, and the only receiver returned at LOS.
- [ ] Waiting for a free worker, like the goroutine on line 13
  > Line 13 takes a worker; line 15 comes after the flight, so these four already hold theirs and never give them back.

## Why does each of those goroutines keep 4 MiB in use?

- [ ] Its stack grew to 4 MiB while it waited
  > A waiting goroutine's stack stays a few kilobytes; the 4 MiB is the pass's buffer, which `make` put on the heap.
- [ ] The collector skips a goroutine's memory until the goroutine runs again
  > The collector looks at every goroutine, running or not, and keeps exactly what they can reach.
- [x] It can still reach its pass's buffer, so no collection can free it
  > A goroutine stuck forever keeps everything it can reach, here its 4 MiB of frames through `rep`.
