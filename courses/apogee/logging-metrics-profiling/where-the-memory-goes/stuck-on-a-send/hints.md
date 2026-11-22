`make` takes a second argument for a channel's room: with room for one result, the goroutine's send finishes even after `runPass` has returned and nobody receives.

A deferred call runs when the goroutine's function returns, however it gets there. On the line right after the worker is taken, defer a small function that receives from `u.workers`, and delete the receive at the goroutine's end.

If the test still counts 34 goroutines, the send on `done` is still waiting: check that the room is in the `make` that creates `done`. If it counts 30 and four workers held, nothing gives the worker back any more: the deferred receive from `u.workers` is missing.
