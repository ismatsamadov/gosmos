Replace the loop over `jobs` with a loop that runs `workers` times, `for range workers`, and start one goroutine with `wg.Go` on each turn.

Inside each worker, range over `jobs` and call `do` for every job you receive, passing any error to `fail`. Every worker ends once `jobs` is closed and empty, so `wg.Wait` still waits for every job.

If the test counts more goroutines than workers while the jobs are held, a goroutine is still started for each job. If only twelve of the eighty jobs run, each worker takes one job instead of ranging over `jobs` until it closes.
