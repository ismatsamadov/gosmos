Make one channel with room for one for each satellite before the week flies, and keep it where every pass of that satellite finds it. In the `select` that reads the call, add a case that sends into it, set only while the pass wants the lock (a nil channel's case never runs), so frames keep coming while it waits.

`max(next, last+1)` numbers each command: keep the satellite's next number beside its lock, touch it only while you hold the lock, and move it on. Let go, with a receive, once the pass has no command left to send or one is refused, and at LOS, whatever is still in flight.

If a satellite refused a command as replayed, its number came from the WELCOME alone, which another pass had already used. If two answers came less than 2 s apart, the second command went up before the first was answered; and if a later pass sent nothing, a lock was never let go.
