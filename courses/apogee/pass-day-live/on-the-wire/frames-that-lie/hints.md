`readFrame` already tells the two apart: a damaged frame comes back with `errBadCRC`, and a length out of range with an error that wraps `errBadLength`. `errors.Is` finds which one came, through the wrapping too.

The loop that reads frames goes on past a damaged one and stops at any other error. On the pass's side, a damaged frame is dropped, and any other error that comes while the pass's context is still running means the line is lost: note it in the pass, and hang up.

If the tests say the rehearsal read fewer TLM frames than Varsk Island relayed, the reading stopped at the bad CRC. If the pass given a length of 8 was flown, or hung up at its LOS, the out-of-range length didn't end the call. An error after the pass's context has ended is the hang-up at LOS, not a lost line.
