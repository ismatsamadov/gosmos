Count the ways `Reading` can return after `Do` succeeds: the status check, a failed decode and the success. Only the last one closes the body today.

A `defer` placed right after the error check of `Do` runs on every one of those returns, so one line covers them all. Then the `Close` call near the end is no longer needed.

If the 404 case still leaves a body open, the `defer` comes after the status check: it has to come before every return it should cover.
