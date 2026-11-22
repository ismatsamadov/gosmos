You don't need to call `TracksDoppler`, only to know whether `r` has it. Which form of type assertion answers that without panicking?

Assert `r` to the `DopplerTracker` interface in the comma-ok form. You can drop the value with `_` and keep only `ok`.

`ok` is `true` for the radios that correct Doppler themselves, and those need no help, so return the opposite of `ok`.
