Each watcher is a `net.Conn`, so it has `SetWriteDeadline`. Call it inside the loop, right before `io.WriteString`.

The instant is `time.Now().Add(w.Timeout)`, worked out again for each write: 2 s, in the tests as at Port Ellery.

If `Broadcast` is still writing after a minute, the sleeping watcher's write has no deadline. If a reading watcher misses a line, one deadline was shared by every write and ran out before its turn.
