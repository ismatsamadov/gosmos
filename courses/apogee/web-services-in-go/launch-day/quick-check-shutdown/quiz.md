## What does `Shutdown` do with a request that is still being answered?

- [ ] Cuts it off, as `Close` does
  > `Close` cuts every connection at once; `Shutdown` exists so that active requests can finish.
- [x] Lets it finish, while refusing new connections at once
  > It closes the listener first, then waits for the requests in flight, or for its context to end.
- [ ] Answers it with `503 Service Unavailable`
  > `Shutdown` writes no answers: a request in flight gets whatever its handler writes.

## Why can the draft above end the process before the last request finishes?

- [x] `Serve` returns as soon as `Shutdown` begins, not when it ends
  > `Serve` returns `ErrServerClosed` the moment the listener closes, so `main` goes on while `Shutdown` is still waiting.
- [ ] `Shutdown` with `context.Background()` gives up at once
  > A background context never ends, so this `Shutdown` would wait as long as it takes, if `main` let it.
- [ ] `log.Println` closes the open connections
  > Logging touches no connection; the process ending does.

## What does `signal.NotifyContext` give `main`?

- [ ] A channel of every signal the process receives
  > That is `signal.Notify`; `NotifyContext` turns a signal into something every function already understands.
- [ ] A function that sends SIGTERM to the old process
  > The deploy tool sends the signal; `main` only listens for it.
- [x] A context that is canceled when one of the signals arrives
  > Code that waits on `ctx.Done()`, such as `run`, then stops when the deploy tool asks.
