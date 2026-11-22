`Start` resumes every delivery the outbox holds, each when it is due. Call it before `serve`, so that Cirrine's waiting event goes out as soon as the new process is up.

`Stop` ends the deliveries in progress and saves each one's attempts to the outbox. Call it only after `stop` has returned: until the shutdown is over, the old process is still at work.

`run` has three returns after serving starts, and each must stop the notifier and join what `Stop` returns to its own error. If the outbox still shows 0 attempts after the first process, `Stop` never ran.
