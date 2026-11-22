`for range ticks` runs its body once per tick and ends when `ticks` closes, which is when `heartbeat` returns `nil`.

In the loop, send a frame of type `TypePing`, then read the next frame with `c.recv()`. Return any error from either, and an error of your own from `fmt.Errorf` when the answer's type isn't `TypePong`.

If the QUERY after the heartbeat gets a PONG, a PING's answer went unread. If the quiet controller's case returns `nil`, no answer is read at all, so no deadline can fire.
