Collect the outcomes in a slice of errors and range over the queue, calling `l.Send` for each command in turn.

`errors.AsType[*NakError](err)` tells a refusal from every other error. A nil error and a refusal both go into the outcomes, and the loop moves on to the next command.

Any other error ends `Drain` at once, returning the outcomes so far and that error. If the wind test stops after the first command, a `*NakError` was taken for the end of the queue.
