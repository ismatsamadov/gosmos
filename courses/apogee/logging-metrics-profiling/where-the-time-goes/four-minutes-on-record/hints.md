`os.Create` makes the file and is the first step that can fail: return its error before anything else runs, so `work` stays unrun when the folder isn't there.

`pprof.StartCPUProfile` takes the file and can fail too, when another profile is running: then close the file and return the error. Otherwise run `work`, call `pprof.StopCPUProfile()`, and return what closing the file returns.

If `cpu.prof` is empty or the test says a profile is still running, `StopCPUProfile` is missing: the profile is written to the file when it stops, so the stop comes before the close.
