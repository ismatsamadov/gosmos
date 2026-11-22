Every `Result` already carries its dish's `Err`, because each goroutine stores it there. Nothing more needs to be recorded.

The loop after `wg.Wait()` throws every result away as soon as one has failed. Without it, `pollAll` returns them all.

If a test reports 0 results, `pollAll` still returns `nil` when a dish fails. Callers such as `watchLine` decide what a failure looks like; `pollAll` only reports it.
