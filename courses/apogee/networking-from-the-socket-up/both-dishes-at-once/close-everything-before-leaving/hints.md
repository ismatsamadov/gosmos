`s.Watchers.CloseAll()` closes every watcher, and each client in `s.Dishes` has a `Close` method.

Put `s.wg.Wait()` last, after every close: it returns once the accept loop and every watcher's goroutine have ended.

If `Shutdown` never returns, it waits before closing, while the watchers' goroutines are still reading. If it returns while a watcher is still being greeted, it doesn't wait for `s.wg`. If a simulator says its connection is still open, a dish client was left out.
