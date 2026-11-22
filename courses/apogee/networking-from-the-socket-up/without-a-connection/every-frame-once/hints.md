`t.Shown` is 0 until the first frame is shown, so it tells you whether there is a last frame to compare with; `t.Last` can't, since 0 is a sequence number too.

When a frame has been shown and `seq` is not more than `t.Last`, count it in `Dropped` and return false before changing anything else.

A frame past `t.Last` jumps over `seq - t.Last - 1` numbers: add them to `Missed`, but only when there was a frame before it. Then set `t.Last`, count it in `Shown` and return true.
