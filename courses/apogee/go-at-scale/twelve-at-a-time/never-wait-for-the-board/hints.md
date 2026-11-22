Loop until the schedule is in: each turn, try `q <- s` in a `select` whose other case is `default`. If the send goes through, return how many you dropped.

When the send can't go, the queue is full: take one schedule out with `<-q` in a second `select` with `default`, add one to `dropped`, and go round again. The second `default` matters when the Board takes a schedule at that very moment.

If the queue holds versions 1 to 4 after offering version 5, your `default` drops the new schedule instead of the oldest one. If the test says `offer` is still waiting, a send or a receive has no `default` beside it.
