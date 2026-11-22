Two loops: the first writes every frame to `c.conn` with `writeFrame`, the second calls `c.recv()` once for each frame sent.

In the second loop, check `recv`'s error before keeping the answer, and return what you have kept so far together with the error.

If a BYE case returns two answers, the second is the empty frame that came back with the error: append only after the error check.
