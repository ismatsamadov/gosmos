Loop over `feeds` and give each one a goroutine with `wg.Go`. Inside it, receive from the feed and send on `out`, each in a `select` beside `<-ctx.Done()`, as `parse` does.

Start one more goroutine that calls `wg.Wait()` and then `close(out)`, and return `out` right after starting it: `merge` itself never waits.

If the case with Tirrin Rock's feed still open fails, `out` closes before every forwarder has returned. A goroutine left after cancel means a forwarder's receive or send has no `ctx.Done()` case.
